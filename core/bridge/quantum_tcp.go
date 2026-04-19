package bridge

import (
	"fmt"
	"io"
	"net"
)

// StartQuantumEgressServer corre en el Nodo Maestro. Recibe conexiones TCP desde los Satelites
// que quieren utilizar la salida a internet global del Maestro.
func StartQuantumEgressServer(listenAddr string) {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Printf("❌ [Egress TCP] Falla al iniciar enrutador maestro en %s: %v\n", listenAddr, err)
		return
	}
	fmt.Printf("🌍 [Egress] Enrutador Global TCP Cuántico activo en %s\n", listenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("❌ [Egress TCP] Error aceptando conexion inter-nodo: %v\n", err)
			continue
		}
		go handleQuantumEgress(conn)
	}
}

func handleQuantumEgress(satConn net.Conn) {
	defer satConn.Close()

	buf := make([]byte, 256)
	
	// 1. Leer "Quantum Handshake" (longitud + direccion objetivo)
	_, err := io.ReadFull(satConn, buf[:1])
	if err != nil {
		return
	}
	
	targetLen := int(buf[0])
	if targetLen == 0 {
		return
	}
	
	_, err = io.ReadFull(satConn, buf[:targetLen])
	if err != nil {
		return
	}
	
	targetAddress := string(buf[:targetLen])
	
	// 2. Resolver y Conectar hacia el internet puro (Google, Steam, etc)
	fmt.Printf("🌍 [Gateway IPv7] Resolviendo y enrutando petición externa hacia: %s\n", targetAddress)
	internetConn, err := net.Dial("tcp", targetAddress)
	if err != nil {
		fmt.Printf("❌ [Gateway IPv7] No se pudo conectar al endpoint externo %s: %v\n", targetAddress, err)
		return
	}
	defer internetConn.Close()

	// 3. Puente Bidireccional
	errc := make(chan error, 2)
	go func() {
		// Lo que viene de internet, va hacia el satelite (Download)
		n, err := io.Copy(satConn, internetConn)
		GlobalTelemetry.BytesSent.Add(n) // Enviamos bytes al túnel IPv7
		GlobalTelemetry.PacketsSent.Add(1) // Contabilizamos 1 stream de paquetes
		errc <- err
	}()
	go func() {
		// Lo que viene del satelite, va hacia internet (Upload)
		n, err := io.Copy(internetConn, satConn)
		GlobalTelemetry.BytesReceived.Add(n)
		GlobalTelemetry.PacketsReceived.Add(1)
		errc <- err
	}()

	<-errc
}
