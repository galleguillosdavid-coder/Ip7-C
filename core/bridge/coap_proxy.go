package bridge

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

// StartCoAPProxy arranca un proxy CoAP UDP en el puerto dado (estándar: 5683).
// Traduce peticiones CoAP de dispositivos con recursos mínimos (sensores a batería)
// a mensajes IEU y viceversa, implementando un subconjunto de CoAP RFC 7252.
func StartCoAPProxy(info *NodeInfo, port int) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("❌ [CoAP] Error resolviendo dirección: %v\n", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("❌ [CoAP] No se pudo abrir puerto %d: %v\n", port, err)
		return
	}
	fmt.Printf("📡 [CoAP] Proxy activo en UDP :%d\n", port)

	go func() {
		buf := make([]byte, 1500) // MTU estándar Ethernet
		for {
			n, clientAddr, err := conn.ReadFromUDP(buf)
			if err != nil {
				continue
			}
			go handleCoAPRequest(conn, clientAddr, buf[:n], info)
		}
	}()
}

// CoAPHeader representa la cabecera fija de 4 bytes de CoAP (RFC 7252)
type CoAPHeader struct {
	Version   uint8  // 2 bits — siempre 1
	Type      uint8  // 2 bits: 0=CON, 1=NON, 2=ACK, 3=RST
	TokenLen  uint8  // 4 bits
	Code      uint8  // 8 bits: clase.detalle (ej. 0.01=GET, 0.02=POST, 2.05=Content)
	MessageID uint16 // 16 bits
}

// CoAPCode valores comunes
const (
	CoAPCodeGET     = 0x01 // 0.01
	CoAPCodePOST    = 0x02 // 0.02
	CoAPCodeContent = 0x45 // 2.05
	CoAPCodeChanged = 0x44 // 2.04
	CoAPCodeBadReq  = 0x80 // 4.00
)

// handleCoAPRequest procesa una petición CoAP y la enruta a la red IEU
func handleCoAPRequest(conn *net.UDPConn, client *net.UDPAddr, data []byte, info *NodeInfo) {
	if len(data) < 4 {
		return
	}

	hdr := CoAPHeader{
		Version:   (data[0] >> 6) & 0x03,
		Type:      (data[0] >> 4) & 0x03,
		TokenLen:  data[0] & 0x0F,
		Code:      data[1],
		MessageID: binary.BigEndian.Uint16(data[2:4]),
	}

	if hdr.Version != 1 {
		sendCoAPResponse(conn, client, hdr, CoAPCodeBadReq, nil)
		return
	}

	// Extraer token y payload
	pos := 4 + int(hdr.TokenLen)
	if pos > len(data) {
		return
	}
	token := data[4 : 4+int(hdr.TokenLen)]

	// Parsear opciones CoAP para extraer Uri-Path
	uriPath := ""
	payload := []byte{}
	for pos < len(data) {
		if data[pos] == 0xFF { // Payload marker
			payload = data[pos+1:]
			break
		}
		optDelta := (data[pos] >> 4) & 0x0F
		optLen := data[pos] & 0x0F
		pos++
		if optDelta == 11 { // Uri-Path option number = 11
			uriPath += "/" + string(data[pos:pos+int(optLen)])
		}
		pos += int(optLen)
	}

	_ = token
	fmt.Printf("📡 [CoAP] %s %s desde %s | %d bytes payload\n",
		coapCodeName(hdr.Code), uriPath, client.String(), len(payload))

	// Enrutamiento por Uri-Path
	switch {
	case strings.HasPrefix(uriPath, "/ieu/send"):
		// POST /ieu/send — encapsular payload como paquete IEU
		if hdr.Code == CoAPCodePOST && len(payload) > 0 {
			if len(payload) < 128 {
				info.Tunnel.SendPriority(payload)
			} else {
				info.Tunnel.SendStandard(payload)
			}
			sendCoAPResponse(conn, client, hdr, CoAPCodeChanged, []byte("OK"))
		} else {
			sendCoAPResponse(conn, client, hdr, CoAPCodeBadReq, []byte("POST con payload requerido"))
		}

	case strings.HasPrefix(uriPath, "/ieu/status"):
		// GET /ieu/status — devolver estado del nodo como JSON-like text
		resp := fmt.Sprintf(`did=%s,latency=%.1fms,version=%s`,
			info.DID, info.Node.Latency, info.Version)
		sendCoAPResponse(conn, client, hdr, CoAPCodeContent, []byte(resp))

	case strings.HasPrefix(uriPath, "/ieu/did"):
		// GET /ieu/did — devolver DID del nodo
		sendCoAPResponse(conn, client, hdr, CoAPCodeContent, []byte(info.DID))

	default:
		sendCoAPResponse(conn, client, hdr, CoAPCodeBadReq, []byte("Ruta CoAP desconocida"))
	}
}

// sendCoAPResponse construye y envía una respuesta CoAP ACK
func sendCoAPResponse(conn *net.UDPConn, client *net.UDPAddr, req CoAPHeader, code uint8, payload []byte) {
	resp := make([]byte, 4)
	resp[0] = (1 << 6) | (2 << 4) // Ver=1, Type=ACK
	resp[1] = code
	binary.BigEndian.PutUint16(resp[2:4], req.MessageID)

	if len(payload) > 0 {
		resp = append(resp, 0xFF) // Payload marker
		resp = append(resp, payload...)
	}
	conn.WriteToUDP(resp, client)
}

func coapCodeName(code uint8) string {
	switch code {
	case CoAPCodeGET:
		return "GET"
	case CoAPCodePOST:
		return "POST"
	default:
		return fmt.Sprintf("0x%02X", code)
	}
}
