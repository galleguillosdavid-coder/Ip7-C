package bridge

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// telemetryWriter envuelve un io.Writer para notificar los bytes escritos en tiempo real
// sin esperar a que la conexión TCP finalice (lo que impedía actualizar el dashboard).
type telemetryWriter struct {
	w  io.Writer
	cb func(int64)
}

func (tw *telemetryWriter) Write(p []byte) (int, error) {
	n, err := tw.w.Write(p)
	if n > 0 && tw.cb != nil {
		tw.cb(int64(n))
	}
	return n, err
}

// Constantes SOCKS5
const (
	socks5Version = 0x05
	addrTypeIPv4  = 0x01
	addrTypeFQDN  = 0x03
	addrTypeIPv6  = 0x04
)

// StartSocks5Server inicia un gateway local en el nodo satélite
// que enruta todo el tráfico a través del Egress TCP Cuántico del Maestro.
func StartSocks5Server(listenAddr string, masterTcpAddr string) {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Printf("❌ [SOCKS5] Falla al iniciar proxy local en %s: %v\n", listenAddr, err)
		return
	}
	fmt.Printf("🌐 [SOCKS5] Proxy Gateway Local activo en %s (Apunta tu navegador aquí)\n", listenAddr)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			fmt.Printf("❌ [SOCKS5] Error aceptando conexion: %v\n", err)
			continue
		}
		go handleSocks5Connection(clientConn, masterTcpAddr)
	}
}

func handleSocks5Connection(client net.Conn, masterTcpAddr string) {
	defer client.Close()
	
	// Incrementamos contadores simulados si es TCP passthrough
	// (Realmente esto se contará en el QuantumTCP pero es bueno sumarlo al nodo local)

	buf := make([]byte, 260)

	// 1. Handshake Inicial SOCKS5
	_, err := io.ReadFull(client, buf[:2])
	if err != nil {
		fmt.Printf("⚠️ [Proxy] Error lectura inicial: %v\n", err)
		return
	}

	// FALLBACK HTTP PROXY: Detectar "CO" de "CONNECT host:port HTTP/1.x"
	// Fix: verificar 2 bytes para evitar falsos positivos con protocolos que empiezan con 'C'
	if buf[0] == 'C' && buf[1] == 'O' {
		// Modo HTTP CONNECT detectado (común cuando los usuarios se confunden y ponen Proxy HTTPS en OS)
		fmt.Println("⚡ [Proxy] Solicitud HTTP CONNECT detectada (Fallback Automático activado)")
		// Buscar el Host en la primera linea "CONNECT host:port HTTP/1.1"
		// Usamos un loop rápido para leer la cabecera
		headerBytes := []byte{buf[0], buf[1]}
		tempBuf := make([]byte, 1)
		for {
			_, err := io.ReadFull(client, tempBuf)
			if err != nil { return }
			headerBytes = append(headerBytes, tempBuf[0])
			l := len(headerBytes)
			if l >= 4 && string(headerBytes[l-4:]) == "\r\n\r\n" {
				break
			}
			if l > 4096 { return } // Demasiado largo
		}
		
		// Parsear el Host
		headerStr := string(headerBytes)
		var targetAddress string
		fmt.Sscanf(headerStr, "CONNECT %s HTTP", &targetAddress)
		
		if targetAddress == "" { return }

		// Conectar al Maestro (IPv7 Quantum Egress Gateway)
		quantumConn, err := net.Dial("tcp", masterTcpAddr)
		if err != nil {
			fmt.Printf("⚠️ [Proxy-HTTP] Error conectando al Gateway Maestro: %v\n", err)
			client.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
			return
		}
		defer quantumConn.Close()

		// Intercambio "Quantum Handshake"
		targetBytes := []byte(targetAddress)
		quantumConn.Write([]byte{byte(len(targetBytes))})
		quantumConn.Write(targetBytes)

		// Responder al cliente HTTP
		client.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

		fmt.Printf("🚇 [HTTP-IPv7] Abriendo tunel hacia Internet via Maestro: %s\n", targetAddress)
		errc := make(chan error, 2)
		go func() {
			_, err := io.Copy(&telemetryWriter{w: quantumConn, cb: func(n int64) { GlobalTelemetry.BytesSent.Add(n) }}, client)
			errc <- err
		}()
		go func() {
			_, err := io.Copy(&telemetryWriter{w: client, cb: func(n int64) { GlobalTelemetry.BytesReceived.Add(n) }}, quantumConn)
			errc <- err
		}()
		<-errc
		return
	}
	
	// Si manda GET crudo
	if buf[0] == 'G' || buf[0] == 'P' {
		fmt.Println("⚠️ [Proxy] ATENCIÓN: El navegador envió una solicitud HTTP Plana.")
		client.Write([]byte("HTTP/1.1 400 Bad Request\r\nContent-Type: text/plain\r\n\r\nPor favor usa HTTPS o configura tu proxy en modo SOCKS5 / CONNECT.\n"))
		return
	}

	if buf[0] != socks5Version {
		fmt.Printf("⚠️ [SOCKS] Protocolo desconocido. Byte0: %v\n", buf[0])
		return
	}

	numMethods := int(buf[1])
	_, err = io.ReadFull(client, buf[:numMethods])
	if err != nil {
		fmt.Printf("⚠️ [SOCKS] Error leyendo methods: %v\n", err)
		return
	}

	// Responder: SOCKS5, No Authentication Required (0x00)
	client.Write([]byte{socks5Version, 0x00})

	// 2. Leer Peticion CONNECT
	_, err = io.ReadFull(client, buf[:4])
	if err != nil || buf[0] != socks5Version || buf[1] != 0x01 { // Solo soportamos CONNECT
		fmt.Printf("⚠️ [SOCKS] Error leyendo CONNECT req: %v. buf0:%v buf1:%v\n", err, buf[0], buf[1])
		return
	}

	addrType := buf[3]
	var destAddr string

	switch addrType {
	case addrTypeIPv4:
		_, err = io.ReadFull(client, buf[:4])
		if err != nil {
			return
		}
		destAddr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])
	case addrTypeFQDN:
		_, err = io.ReadFull(client, buf[:1])
		if err != nil {
			return
		}
		domainLen := int(buf[0])
		_, err = io.ReadFull(client, buf[:domainLen])
		if err != nil {
			return
		}
		destAddr = string(buf[:domainLen])
	case addrTypeIPv6:
		// Simplificado, omitimos IPv6 raw por ahora
		fmt.Printf("⚠️ [SOCKS] IPv6 not supported\n")
		return
	default:
		fmt.Printf("⚠️ [SOCKS] Unsupported address type: %v\n", addrType)
		return
	}

	// Leer Puerto Destino
	_, err = io.ReadFull(client, buf[:2])
	if err != nil {
		fmt.Printf("⚠️ [SOCKS] Error leyendo puerto: %v\n", err)
		return
	}
	destPort := binary.BigEndian.Uint16(buf[:2])
	targetAddress := fmt.Sprintf("%s:%d", destAddr, destPort)

	// 3. Conectar al Maestro (IPv7 Quantum Egress Gateway)
	// Para mayor velocidad y estabilidad en descargas pesadas saltamos UDP interno
	// y usamos un canal TCP secundario (que el Master escucha) y enviamos un "Protocol Header IPv7" para validar.
	
	quantumConn, err := net.Dial("tcp", masterTcpAddr)
	if err != nil {
		fmt.Printf("⚠️ [SOCKS5] Error conectando al Gateway Maestro IPv7: %v\n", err)
		// Responder fallo al cliente SOCKS5 (Host uncreachable)
		client.Write([]byte{socks5Version, 0x04, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
		return
	}
	defer quantumConn.Close()

	// Intercambio "Quantum Handshake": Enviamos qué destino queremos abrir
	// Formato: [Longitud 1 byte][Target Address]
	targetBytes := []byte(targetAddress)
	if len(targetBytes) > 255 {
		return // Fallback de seguridad
	}
	quantumConn.Write([]byte{byte(len(targetBytes))})
	quantumConn.Write(targetBytes)

	// Responder OK al cliente local SOCKS5 Chrome/Windows
	// Indicamos exito (0x00) y rellenamos address generico
	client.Write([]byte{socks5Version, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})

	// 4. Iniciar Copy Bidireccional
	fmt.Printf("🚇 [SOCKS-IPv7] Abriendo tunel hacia Internet via Maestro: %s\n", targetAddress)
	
	errc := make(chan error, 2)
	go func() {
		_, err := io.Copy(&telemetryWriter{w: quantumConn, cb: func(n int64) { GlobalTelemetry.BytesSent.Add(n) }}, client)
		errc <- err
	}()
	go func() {
		_, err := io.Copy(&telemetryWriter{w: client, cb: func(n int64) { GlobalTelemetry.BytesReceived.Add(n) }}, quantumConn)
		errc <- err
	}()

	<-errc
}
