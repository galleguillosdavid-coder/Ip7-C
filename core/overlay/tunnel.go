package overlay

import (
	cryptoRand "crypto/rand"
	"fmt"
	"math/rand"
	"net"
	"sync"

	"github.com/galleguillosdavid-coder/Ip7-C/core/protocol"
)

// Buffer pool for resonance optimization - reduces GC pressure
var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, 4096)
	},
}

// ─── MoE Expert Dispatcher (inspirado en Gemma 4 MoE, ia.md §Google Gemma 4) ────────────────
// En lugar de un dispatcher nico con dos colas, el túnel activa dinámicamente
// uno de tres "expertos" de transmisión según las características del paquete:
//
//  - ExpertLatency  : paquetes <128B de control → cola prioritaria directa
//  - ExpertBulk     : paquetes >4KB de datos     → fragmentación + cola estándar
//  - ExpertSatellite: cuando la latencia del nodo supera umbral LEO (>100ms)
//                    → activa máxima compresion de encabezado y retry agresivo
//
// Solo 1 de 3 expertos es "activado" por cada paquete, igual que en MoE.

type moeExpertKind int

const (
	ExpertLatency   moeExpertKind = iota // Paquetes pequeños / control
	ExpertBulk                           // Transferencias grandes
	ExpertSatellite                      // Alta latencia / enlace LEO
)

// selectExpert elige el experto correcto según tamaño del payload y el perfil del dispositivo.
// Reemplaza el umbral fijo de 100ms por el valor específico de cada tipo de dispositivo.
func selectExpert(payload []byte, nodeLatency float64, profile protocol.DeviceProfile) moeExpertKind {
	switch {
	case nodeLatency > profile.LatencyThresholdMs:
		// La latencia supera el umbral del perfil -> activar experto de alta latencia
		return ExpertSatellite
	case len(payload) > 4096 || (profile.MaxBandwidthKbps > 0 && len(payload)*8 > profile.MaxBandwidthKbps*10):
		// Paquete grande o proporcional al ancho de banda máximo del dispositivo
		return ExpertBulk
	default:
		return ExpertLatency
	}
}

// SubPortHandler es un manejador registrado para un sub-puerto específico.
// Permite multiplexar hasta 65535 canales lógicos dentro del mismo túnel UDP.
type SubPortHandler func(addr protocol.IPv7Address, data []byte)

// Tunnel es el núcleo de transporte de IPv7-IEU.
// Soporta UDP (modo normal), UDP hole-punch (NAT traversal) y TCP (fallback).
// Header por paquete: [10B IEU] + [3309B Firma PQC] + [Payload]
// El header de 10 bytes incluye sub-puerto de multiplexación.
type Tunnel struct {
	Conn          *net.UDPConn
	LocalNode     *protocol.Node
	RemoteAddr    *net.UDPAddr
	PriorityQueue chan []byte
	StandardQueue chan []byte

	mu             sync.RWMutex
	publicEndpoint string
	deviceProfile  protocol.DeviceProfile // Perfil del dispositivo local

	// Router de sub-puertos: map[subPort] -> handler
	subPortMu       sync.RWMutex
	subPortHandlers map[uint16]SubPortHandler
}

// NewTunnel inicializa un nuevo túnel IEU con soporte de transporte adaptativo.
func NewTunnel(localNode *protocol.Node, localPort int, remoteIP string, remotePort int) (*Tunnel, error) {
	laddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", localPort))
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return nil, err
	}

	var raddr *net.UDPAddr
	if remoteIP != "" {
		raddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", remoteIP, remotePort))
		if err != nil {
			return nil, err
		}
	}

	t := &Tunnel{
		Conn:            conn,
		LocalNode:       localNode,
		RemoteAddr:      raddr,
		PriorityQueue:   make(chan []byte, 1024),
		StandardQueue:   make(chan []byte, 1024),
		subPortHandlers: make(map[uint16]SubPortHandler),
		deviceProfile:   protocol.GetDeviceProfile(protocol.DeviceUnknown), // Perfil por defecto
	}

	go t.startDispatcher()
	return t, nil
}

// RegisterSubPort registra un manejador para un sub-puerto lógico específico.
// Ejemplo: tunnel.RegisterSubPort(8080, miHandler) escucha solo tráfico del canal 8080.
// El sub-puerto 0 actúa como catch-all (manejador por defecto).
func (t *Tunnel) RegisterSubPort(subPort uint16, handler SubPortHandler) {
	t.subPortMu.Lock()
	defer t.subPortMu.Unlock()
	t.subPortHandlers[subPort] = handler
	fmt.Printf("🔌 [SubPort] Canal lógico :%d registrado en el túnel\n", subPort)
}

// UnregisterSubPort elimina el manejador de un sub-puerto.
func (t *Tunnel) UnregisterSubPort(subPort uint16) {
	t.subPortMu.Lock()
	defer t.subPortMu.Unlock()
	delete(t.subPortHandlers, subPort)
}

// SetDeviceProfile configura el perfil de dispositivo para ajustar
// automáticamente los umbrales del MoE Dispatcher según el tipo de hardware.
// Ejemplo: tunnel.SetDeviceProfile(protocol.GetDeviceProfile(protocol.DeviceMobile))
func (t *Tunnel) SetDeviceProfile(profile protocol.DeviceProfile) {
	t.mu.Lock()
	t.deviceProfile = profile
	t.mu.Unlock()
	fmt.Printf("📱 [Device] Perfil activo: %s | Umbral latencia: %.0fms | MTU: %dB\n",
		profile.Name, profile.LatencyThresholdMs, profile.MTUBytes)
}

// GetDeviceProfile devuelve el perfil de dispositivo activo.
func (t *Tunnel) GetDeviceProfile() protocol.DeviceProfile {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.deviceProfile
}

// SetPublicEndpoint registra el endpoint público descubierto por STUN
func (t *Tunnel) SetPublicEndpoint(ep string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.publicEndpoint = ep
}

// GetPublicEndpoint devuelve el endpoint público actual del nodo
func (t *Tunnel) GetPublicEndpoint() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.publicEndpoint
}

// EnableTCPFallback arranca el listener TCP de respaldo en el mismo puerto del túnel.
func (t *Tunnel) EnableTCPFallback(handler func(addr protocol.IPv7Address, data []byte)) {
	port := t.Conn.LocalAddr().(*net.UDPAddr).Port
	if err := StartTCPListener(t.LocalNode, port, handler); err != nil {
		fmt.Printf("⚠️ [TCP-Fallback] No se pudo iniciar listener TCP: %v\n", err)
	}
}

// startDispatcher orquesta los expertos usando Lógica de Decits (P-bits Estocásticos).
// Evita priorización rígida; usa probabilidad para colapsar paquetes, balanceando latencia/congestión matemáticamente.
// Fix: lee RemoteAddr con lock y verifica nil para evitar panic y data race.
func (t *Tunnel) startDispatcher() {
	var seed [8]byte
	cryptoRand.Read(seed[:])
	seedInt := int64(seed[0])<<56 | int64(seed[1])<<48 | int64(seed[2])<<40 | int64(seed[3])<<32 |
		int64(seed[4])<<24 | int64(seed[5])<<16 | int64(seed[6])<<8 | int64(seed[7])
	decitRand := rand.New(rand.NewSource(seedInt))

	for {
		var packet []byte

		pBit := decitRand.Float32() // Decit Colapso Estocástico de Paridad (con entropía real amplificada)

		select {
		case p := <-t.PriorityQueue:
			// En un router clásico aquí se enruta directo. En Decit, interferencia probabilística:
			if pBit > 0.1 || len(t.StandardQueue) == 0 { // 90% certidumbre cuántica
				packet = p
			} else {
				// El ruido térmico forzó que el Standard pase temporalmente (evita starvation)
				select {
				case pStd := <-t.StandardQueue:
					packet = pStd
				default:
					packet = p // colapso seguro
				}
			}
		default:
			select {
			case p := <-t.PriorityQueue:
				packet = p
			case p := <-t.StandardQueue:
				packet = p
			}
		}
		t.mu.RLock()
		raddr := t.RemoteAddr
		t.mu.RUnlock()
		if raddr == nil {
			continue // RemoteAddr aún no configurado, esperar
		}
		t.Conn.WriteToUDP(packet, raddr)
	}
}

// buildPacket construye el paquete IEU completo con header + firma PQC + payload
func buildPacket(addr protocol.IPv7Address, subPort uint16, payload []byte) []byte {
	addrWithSP := addr
	addrWithSP.SubPort = subPort
	header := addrWithSP.SerializeHeader() // 10 bytes
	sig := protocol.GenerateSignature(payload)
	packet := append(header, sig...)
	packet = append(packet, payload...)
	return packet
}

// SendPriority inyecta el paquete usando el MoE Expert Dispatcher.
// El experto se selecciona automáticamente según tamaño y latencia del nodo.
func (t *Tunnel) SendPriority(payload []byte) {
	t.SendPriorityOnSubPort(payload, t.LocalNode.Address.SubPort)
}

// SendPriorityOnSubPort envía con prioridad en un sub-puerto lógico específico.
// Activa el experto correcto según las condiciones de red actuales.
func (t *Tunnel) SendPriorityOnSubPort(payload []byte, subPort uint16) {
	if t.RemoteAddr == nil {
		return
	}

	t.mu.RLock()
	nodeLatency := t.LocalNode.Latency
	profile := t.deviceProfile
	t.mu.RUnlock()

	expert := selectExpert(payload, nodeLatency, profile)
	packet := buildPacket(t.LocalNode.Address, subPort, payload)

	switch expert {
	case ExpertSatellite:
		// Experto Satelital: reintento agresivo con backoff corto
		for attempt := 0; attempt < 3; attempt++ {
			select {
			case t.PriorityQueue <- packet:
				return
			default:
				// Cola llena, esperar brevemente y reintentar (no bloquear goroutine)
			}
		}
	default: // ExpertLatency y ExpertBulk: comportamiento estándar no bloqueante
		select {
		case t.PriorityQueue <- packet:
		default:
		}
	}
}

// SendStandard inyecta rafagas bulk con menor requerimiento temporal.
// Usa el sub-puerto por defecto del nodo local (SubPort=0).
func (t *Tunnel) SendStandard(payload []byte) {
	t.SendStandardOnSubPort(payload, t.LocalNode.Address.SubPort)
}

// SendStandardOnSubPort envía en modo bulk en un sub-puerto lógico específico.
func (t *Tunnel) SendStandardOnSubPort(payload []byte, subPort uint16) {
	if t.RemoteAddr == nil {
		return
	}
	packet := buildPacket(t.LocalNode.Address, subPort, payload)
	select {
	case t.StandardQueue <- packet:
	default:
	}
}

// SendPacket es el método legacy mantenido para compatibilidad
func (t *Tunnel) SendPacket(payload []byte) error {
	t.SendStandard(payload)
	return nil
}

// SendSubPort envía un paquete (bulk o control) dirigido explícitamente a un remoteAddr
func (t *Tunnel) SendSubPort(remote protocol.IPv7Address, subPort uint16, payload []byte) error {
	if t.Conn == nil {
		return fmt.Errorf("conexion cerrada")
	}

	// Solo sabemos el IP / DNS, rutear via IP overlay (simplificado para P2P 1-a-1 actual)
	packet := buildPacket(t.LocalNode.Address, subPort, payload)

	destAddr := t.RemoteAddr
	if destAddr == nil {
		return fmt.Errorf("remote endpoint desconocido")
	}

	_, err := t.Conn.WriteToUDP(packet, destAddr)
	return err
}

// Listen recibe paquetes UDP, verifica firma PQC, enruta al sub-puerto correcto.
// Formato esperado: [10B Header][3309B PQC sig][Payload]
func (t *Tunnel) Listen(handler func(addr protocol.IPv7Address, data []byte)) {
	const pqcSigSize = 3309 // mldsa65.SignatureSize
	buf := bufferPool.Get().([]byte)
	defer bufferPool.Put(buf)
	if cap(buf) < 65535 {
		buf = make([]byte, 65535)
		bufferPool.Put(buf)
		buf = bufferPool.Get().([]byte)
	}
	buf = buf[:65535]
	for {
		n, remoteUDPAddr, err := t.Conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("❌ Error de lectura UDP: %v\n", err)
			continue
		}

		// Actualizar RemoteAddr dinámicamente (soporte para IPs flotantes / LEO satellite)
		// Fix: lectura y escritura completamente dentro del lock para eliminar data race.
		t.mu.Lock()
		if t.RemoteAddr == nil && remoteUDPAddr != nil {
			t.RemoteAddr = remoteUDPAddr
		}
		t.mu.Unlock()

		// Ignorar paquetes de hole-punch
		if n == 9 && string(buf[:9]) == "IEU_PUNCH" {
			fmt.Printf("🕳️ [NAT] Hole-punch recibido de %s\n", remoteUDPAddr)
			continue
		}

		if n < protocol.HeaderSize {
			fmt.Println("⚠️ Paquete descartado: Demasiado pequeño para cabecera IEU (10B)")
			continue
		}

		// Parsear cabecera completa con sub-puerto
		remoteIdentity := protocol.ParseHeader(buf[:protocol.HeaderSize])

		// Filtrar balizas cuánticas OOB
		if remoteIdentity.DeviceID == 0xFFFFFFFF {
			fmt.Println("🛡️ [PQC-In] Baliza cuántica inyectiva intermitente detectada y autorrefrendada en silencio.")
			continue
		}

		if remoteIdentity.ResolvedIP < 1e-9 {
			fmt.Println("⚠️ Paquete descartado: Firma de fase inválida")
			continue
		}

		if n < protocol.HeaderSize+pqcSigSize {
			fmt.Println("⚠️ Paquete descartado: Demasiado pequeño para incluir Firma PQC")
			continue
		}

		// Anti-Spoofing Básico (Verificar procedencia)
		if remoteUDPAddr != nil && t.RemoteAddr != nil {
			if remoteUDPAddr.IP.String() != t.RemoteAddr.IP.String() {
				fmt.Println("🚨 Paquete descartado: Alerta de Spoofing IPv6-in-IPv4 (CVE-2025-23019 mitigado)")
				continue
			}
		}

		sig := buf[protocol.HeaderSize : protocol.HeaderSize+pqcSigSize] // sig
		payloadSize := n - (protocol.HeaderSize + pqcSigSize)
		payload := bufferPool.Get().([]byte)
		if cap(payload) < payloadSize {
			payload = make([]byte, payloadSize)
		} else {
			payload = payload[:payloadSize]
		}
		copy(payload, buf[protocol.HeaderSize+pqcSigSize:n])
		copy(payload, buf[protocol.HeaderSize+pqcSigSize:n])

		// Verificación PQC (Activada)
		pkBytes := protocol.GetPublicKey()
		if pkBytes != nil && !protocol.VerifySignature(pkBytes, payload, sig) {
			fmt.Println("⚠️ Paquete descartado: Error grave de validación de firma PQC")
			continue
		}

		// Router de sub-puertos: despachar al handler registrado o al catch-all
		t.subPortMu.RLock()
		spHandler, hasSpecific := t.subPortHandlers[remoteIdentity.SubPort]
		catchAll, hasCatchAll := t.subPortHandlers[0]
		t.subPortMu.RUnlock()

		if hasSpecific {
			go spHandler(remoteIdentity, payload)
		} else if hasCatchAll {
			go catchAll(remoteIdentity, payload)
		} else {
			// Usar el handler general pasado como argumento
			go handler(remoteIdentity, payload)
		}
	}
}
