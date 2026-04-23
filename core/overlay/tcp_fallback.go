package overlay

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/galleguillosdavid-coder/Ip7-C/core/protocol"
)

// TCPFallbackPorts lista de puertos de fallback en orden de prioridad
// 443 = camouflage como HTTPS (pasa casi cualquier firewall)
var TCPFallbackPorts = []int{443, 8443, 80, 8080}

// TCPSession representa una conexión TCP de respaldo para túneles bloqueados
type TCPSession struct {
	conn      net.Conn
	localNode *protocol.Node
}

// StartTCPListener arranca un servidor TCP en el puerto dado para aceptar
// conexiones de nodos que no puedan llegar via UDP.
func StartTCPListener(localNode *protocol.Node, port int, handler func(addr protocol.IPv7Address, data []byte)) error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("error iniciando listener TCP en puerto %d: %v", port, err)
	}
	fmt.Printf("🔁 [TCP-Fallback] Listener activo en :%d\n", port)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				continue
			}
			go handleTCPSession(conn, localNode, handler)
		}
	}()
	return nil
}

// handleTCPSession procesa paquetes IEU encapsulados en una sesión TCP.
// Formato de framing: [4 bytes longitud big-endian][N bytes payload IEU]
func handleTCPSession(conn net.Conn, localNode *protocol.Node, handler func(addr protocol.IPv7Address, data []byte)) {
	defer conn.Close()

	lenBuf := make([]byte, 4)
	for {
		// Renovar deadline en cada frame: sesiones largas no se cortan silenciosamente
		conn.SetDeadline(time.Now().Add(10 * time.Minute))

		if _, err := readFull(conn, lenBuf); err != nil {
			return
		}
		frameLen := binary.BigEndian.Uint32(lenBuf)
		if frameLen > 65535 {
			return // MTU excedido, conexión inválida
		}

		frame := make([]byte, frameLen)
		if _, err := readFull(conn, frame); err != nil {
			return
		}

		if len(frame) < 8 {
			continue
		}

		remoteIdentity := protocol.ParseHeader(frame[:8])
		if remoteIdentity.ResolvedIP < 1e-9 {
			continue // firma inválida
		}
		handler(remoteIdentity, frame[8:])
	}
}

// ConnectTCPFallback intenta conectar a un peer via TCP cuando UDP falla.
// Prueba los puertos de camouflage en orden hasta encontrar uno accesible.
func ConnectTCPFallback(remoteHost string) (*TCPSession, error) {
	for _, port := range TCPFallbackPorts {
		addr := net.JoinHostPort(remoteHost, strconv.Itoa(port))
		conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
		if err != nil {
			continue
		}
		fmt.Printf("✅ [TCP-Fallback] Conectado via TCP en %s\n", addr)
		return &TCPSession{conn: conn}, nil
	}
	return nil, fmt.Errorf("no se pudo conectar a %s via TCP en ningún puerto de fallback", remoteHost)
}

// Send encapsula y envía un frame IEU sobre la sesión TCP con framing de longitud
func (s *TCPSession) Send(localNode *protocol.Node, payload []byte) error {
	header := localNode.Address.SerializeHeader()
	frame := append(header, payload...)
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(frame)))
	full := append(lenBuf, frame...)
	_, err := s.conn.Write(full)
	return err
}

// Close cierra la sesión TCP de fallback
func (s *TCPSession) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

// readFull lee exactamente len(buf) bytes de la conexión
func readFull(conn net.Conn, buf []byte) (int, error) {
	total := 0
	for total < len(buf) {
		n, err := conn.Read(buf[total:])
		total += n
		if err != nil {
			return total, err
		}
	}
	return total, nil
}
