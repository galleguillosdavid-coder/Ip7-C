package p2p

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	KeySize = 20
)

type NodeID [KeySize]byte

func HashString(data string) NodeID {
	return sha1.Sum([]byte(data))
}

type Peer struct {
	ID   NodeID
	Addr *net.UDPAddr
}

// MicroDHT implementa un anillo Web3 hiper-liviano estilo Kademlia para IPv7.
// Incorpora Task Budget (inspirado en Claude 4.7, ia.md §Anthropic Claude 4.7):
// cada nodo tiene un presupuesto de broadcasts por ventana temporal para evitar
// saturar enlaces satelitales de alta latencia (Starlink ~20ms RTT pero limitado ancho de banda).
type MicroDHT struct {
	LocalID NodeID
	Conn    *net.UDPConn

	dbMu        sync.RWMutex
	Buckets     [KeySize * 8][]*Peer
	bucketLocks [KeySize * 8]sync.RWMutex
	DB          map[NodeID]string // Almacén en Memoria Corta: Hash(DID) -> IP:Port

	// Task Budget: limita RPC de salida para no saturar enlaces LEO/satelitales
	budgetMu        sync.Mutex
	budgetUsed      int       // RPCs enviados en la ventana actual
	budgetMax       int       // Máximo RPCs por ventana
	budgetWindowEnd time.Time // Fin de la ventana de tiempo actual
	budgetWindow    time.Duration
}

type RPCMessage struct {
	Type   string
	Sender NodeID
	Key    NodeID
	Value  string
}

// NewMicroDHT arranca un nodo rastreador satelital P2P independiente.
// Acepta un presupuesto de broadcasts máximo por minuto (0 = sin límite).
// Recomendado: 60 para enlaces de fibra, 20 para Starlink LEO.
func NewMicroDHT(localID string, port int) (*MicroDHT, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	dht := &MicroDHT{
		LocalID:         HashString(localID),
		Conn:            conn,
		DB:              make(map[NodeID]string),
		budgetMax:       60, // 60 RPCs por minuto por defecto
		budgetWindow:    time.Minute,
		budgetWindowEnd: time.Now().Add(time.Minute),
	}

	go dht.listen()
	return dht, nil
}

func (dht *MicroDHT) listen() {
	buf := make([]byte, 2048)
	for {
		n, addr, err := dht.Conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		var msg RPCMessage
		if err := json.Unmarshal(buf[:n], &msg); err != nil {
			continue
		}
		dht.handleRPC(msg, addr)
	}
}

func (dht *MicroDHT) handleRPC(msg RPCMessage, addr *net.UDPAddr) {
	dht.updateBucket(msg.Sender, addr)

	switch msg.Type {
	case "STORE":
		dht.dbMu.Lock()
		dht.DB[msg.Key] = msg.Value
		dht.dbMu.Unlock()
	case "FIND_VALUE":
		dht.dbMu.RLock()
		val, exists := dht.DB[msg.Key]
		dht.dbMu.RUnlock()
		if exists {
			reply := RPCMessage{Type: "FOUND", Sender: dht.LocalID, Key: msg.Key, Value: val}
			dht.sendRPC(addr, reply)
		}
	case "FOUND":
		dht.dbMu.Lock()
		dht.DB[msg.Key] = msg.Value
		dht.dbMu.Unlock()
	}
}

// bucketIndex calcula el índice de bucket Kademlia usando distancia XOR.
// Retorna la posición del bit más significativo en XOR(localID, peerID), rango 0..159.
// Esto implementa correctamente la estructura de k-buckets del paper Kademlia original.
func (dht *MicroDHT) bucketIndex(id NodeID) int {
	for i := 0; i < KeySize; i++ {
		xor := dht.LocalID[i] ^ id[i]
		if xor != 0 {
			for bit := 7; bit >= 0; bit-- {
				if xor>>uint(bit) != 0 {
					return i*8 + (7 - bit)
				}
			}
		}
	}
	return KeySize*8 - 1
}

func (dht *MicroDHT) updateBucket(id NodeID, addr *net.UDPAddr) {
	if id == dht.LocalID {
		return // No agregar a sí mismo
	}
	idx := dht.bucketIndex(id)
	dht.bucketLocks[idx].Lock()
	defer dht.bucketLocks[idx].Unlock()
	bucket := dht.Buckets[idx]
	for _, p := range bucket {
		if p.ID == id {
			p.Addr = addr
			return
		}
	}
	if len(bucket) < 20 { // k=20: tamaño estándar de k-bucket Kademlia
		dht.Buckets[idx] = append(bucket, &Peer{ID: id, Addr: addr})
	}
}

// sendRPC envía un mensaje RPC respetando el Task Budget.
// Si el presupuesto de la ventana actual está agotado, el mensaje se descarta
// para proteger el ancho de banda en enlaces satelitales.
func (dht *MicroDHT) sendRPC(addr *net.UDPAddr, msg RPCMessage) {
	dht.budgetMu.Lock()
	now := time.Now()
	if now.After(dht.budgetWindowEnd) {
		// Nueva ventana: reiniciar presupuesto
		dht.budgetUsed = 0
		dht.budgetWindowEnd = now.Add(dht.budgetWindow)
	}
	if dht.budgetMax > 0 && dht.budgetUsed >= dht.budgetMax {
		dht.budgetMu.Unlock()
		// Presupuesto agotado: silenciosamente descartar el RPC
		// (equivalente al comportamiento de Task Budget en Claude 4.7)
		return
	}
	dht.budgetUsed++
	dht.budgetMu.Unlock()

	b, _ := json.Marshal(msg)
	dht.Conn.WriteToUDP(b, addr)
}

// SetBudget ajusta el presupuesto de RPCs por ventana en tiempo de ejecución.
// Útil para adaptar dinámicamente según la calidad del enlace:
//   - Enlace de fibra: SetBudget(120, time.Minute)
//   - Starlink LEO:    SetBudget(20, time.Minute)
//   - Sin límite:     SetBudget(0, time.Minute)
func (dht *MicroDHT) SetBudget(maxRPCsPerWindow int, window time.Duration) {
	dht.budgetMu.Lock()
	dht.budgetMax = maxRPCsPerWindow
	dht.budgetWindow = window
	dht.budgetWindowEnd = time.Now().Add(window)
	dht.budgetMu.Unlock()
}

// Announce propaga el Identificador Descentralizado (DID) W3C a la telaraña global
func (dht *MicroDHT) Announce(key string, value string) {
	hashKey := HashString(key)
	dht.dbMu.Lock()
	dht.DB[hashKey] = value
	dht.dbMu.Unlock()

	dht.bucketLocks[0].RLock()
	peers := append([]*Peer(nil), dht.Buckets[0]...)
	dht.bucketLocks[0].RUnlock()

	msg := RPCMessage{Type: "STORE", Sender: dht.LocalID, Key: hashKey, Value: value}
	for _, peer := range peers {
		dht.sendRPC(peer.Addr, msg)
	}
}

// Resolve busca en la red de manera asincrónica tolerando altas latencias LEO/GEO
func (dht *MicroDHT) Resolve(key string) string {
	hashKey := HashString(key)
	dht.dbMu.RLock()
	val, ok := dht.DB[hashKey]
	dht.dbMu.RUnlock()

	dht.bucketLocks[0].RLock()
	peers := append([]*Peer(nil), dht.Buckets[0]...)
	dht.bucketLocks[0].RUnlock()

	if ok {
		return val
	}
	if len(peers) == 0 {
		return ""
	}

	msg := RPCMessage{Type: "FIND_VALUE", Sender: dht.LocalID, Key: hashKey}
	for _, peer := range peers {
		dht.sendRPC(peer.Addr, msg)
	}

	// Sondeo P2P tolerante usando concurrencia y canales (optimizado para no bloquear threads pasivos)
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	timeout := time.After(3 * time.Second) // 3 segundos max timeout Satelital LEO

	for {
		select {
		case <-timeout:
			return ""
		case <-ticker.C:
			dht.dbMu.RLock()
			val, ok := dht.DB[hashKey]
			dht.dbMu.RUnlock()
			if ok {
				return val
			}
		}
	}
}

// GetPeerList devuelve la lista de peers conocidos como strings "IP:Puerto"
// para ser consumida por la REST API y el WoT descriptor.
func (dht *MicroDHT) GetPeerList() []string {
	var peers []string
	for idx := 0; idx < KeySize*8; idx++ {
		dht.bucketLocks[idx].RLock()
		for _, p := range dht.Buckets[idx] {
			if p.Addr != nil {
				peers = append(peers, p.Addr.String())
			}
		}
		dht.bucketLocks[idx].RUnlock()
	}
	return peers
}

// AddBootstrap anula la centralización conectando a un enjambre estocástico
func (dht *MicroDHT) AddBootstrap(addrStr string) {
	addr, err := net.ResolveUDPAddr("udp", addrStr)
	if err == nil {
		msg := RPCMessage{Type: "PING", Sender: dht.LocalID}
		dht.sendRPC(addr, msg)
	}
}
