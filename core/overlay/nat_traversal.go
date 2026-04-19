package overlay

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// STUNServer lista de servidores STUN públicos (Google + Cloudflare)
var stunServers = []string{
	"stun.l.google.com:19302",
	"stun1.l.google.com:19302",
	"stun.cloudflare.com:3478",
}

// DiscoverPublicEndpoint consulta un servidor STUN para descubrir la IP:Puerto público del nodo.
// Implementa el protocolo STUN RFC 5389 con Binding Request/Response.
func DiscoverPublicEndpoint() (string, error) {
	for _, server := range stunServers {
		ep, err := querySTUN(server)
		if err == nil {
			return ep, nil
		}
	}
	return "", fmt.Errorf("todos los servidores STUN fallaron")
}

// querySTUN envía un STUN Binding Request y parsea el XOR-MAPPED-ADDRESS de la respuesta
func querySTUN(serverAddr string) (string, error) {
	conn, err := net.DialTimeout("udp", serverAddr, 3*time.Second)
	if err != nil {
		return "", fmt.Errorf("no se pudo conectar a %s: %v", serverAddr, err)
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	// STUN Binding Request (20 bytes)
	req := make([]byte, 20)
	// Tipo: 0x0001 = Binding Request
	binary.BigEndian.PutUint16(req[0:2], 0x0001)
	// Longitud del mensaje (sin cabecera): 0 atributos
	binary.BigEndian.PutUint16(req[2:4], 0)
	// Magic Cookie: 0x2112A442 (RFC 5389)
	binary.BigEndian.PutUint32(req[4:8], 0x2112A442)
	// Transaction ID: 12 bytes aleatorios (simplificado con timestamp)
	tid := uint64(time.Now().UnixNano())
	binary.BigEndian.PutUint64(req[8:16], tid)
	binary.BigEndian.PutUint32(req[16:20], 0xDEADBEEF)

	if _, err := conn.Write(req); err != nil {
		return "", fmt.Errorf("error enviando STUN request: %v", err)
	}

	resp := make([]byte, 512)
	n, err := conn.Read(resp)
	if err != nil {
		return "", fmt.Errorf("error leyendo STUN response: %v", err)
	}
	if n < 20 {
		return "", fmt.Errorf("respuesta STUN demasiado corta: %d bytes", n)
	}

	// Verificar tipo de respuesta: 0x0101 = Binding Success Response
	respType := binary.BigEndian.Uint16(resp[0:2])
	if respType != 0x0101 {
		return "", fmt.Errorf("respuesta STUN inesperada, tipo: 0x%04X", respType)
	}

	// Parsear atributos STUN buscando XOR-MAPPED-ADDRESS (0x0020) o MAPPED-ADDRESS (0x0001)
	pos := 20
	for pos+4 <= n {
		attrType := binary.BigEndian.Uint16(resp[pos : pos+2])
		attrLen := int(binary.BigEndian.Uint16(resp[pos+2 : pos+4]))
		pos += 4

		if pos+attrLen > n {
			break
		}

		switch attrType {
		case 0x0020: // XOR-MAPPED-ADDRESS
			if attrLen < 8 {
				break
			}
			family := resp[pos+1]
			if family != 0x01 { // IPv4
				break
			}
			xorPort := binary.BigEndian.Uint16(resp[pos+2:pos+4]) ^ 0x2112
			xorIP := binary.BigEndian.Uint32(resp[pos+4:pos+8]) ^ 0x2112A442
			ip := net.IP([]byte{
				byte(xorIP >> 24),
				byte(xorIP >> 16),
				byte(xorIP >> 8),
				byte(xorIP),
			})
			return fmt.Sprintf("%s:%d", ip.String(), xorPort), nil

		case 0x0001: // MAPPED-ADDRESS (fallback RFC 3489)
			if attrLen < 8 {
				break
			}
			port := binary.BigEndian.Uint16(resp[pos+2 : pos+4])
			ip := net.IP(resp[pos+4 : pos+8])
			return fmt.Sprintf("%s:%d", ip.String(), port), nil
		}

		// Avanzar al siguiente atributo (con padding a 4 bytes)
		pad := (4 - (attrLen % 4)) % 4
		pos += attrLen + pad
	}

	return "", fmt.Errorf("XOR-MAPPED-ADDRESS no encontrado en respuesta STUN")
}

// HolePunch intenta perforar un NAT enviando paquetes UDP al endpoint remoto
// antes de que el otro lado tenga la dirección, creando una entrada en el NAT table.
// Ambos lados deben ejecutar esto simultáneamente (coordinado via MicroDHT).
func HolePunch(localConn *net.UDPConn, remoteEndpoint string, attempts int) error {
	raddr, err := net.ResolveUDPAddr("udp", remoteEndpoint)
	if err != nil {
		return fmt.Errorf("dirección de hole-punch inválida: %v", err)
	}

	punch := []byte("IEU_PUNCH") // marcador de hole punch
	for i := 0; i < attempts; i++ {
		localConn.WriteToUDP(punch, raddr)
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Printf("🕳️ [NAT] Hole-punch enviado a %s (%d intentos)\n", remoteEndpoint, attempts)
	return nil
}
