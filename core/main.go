package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/galleguillosdavid-coder/Ip7-C/core/adapter"
	"github.com/galleguillosdavid-coder/Ip7-C/core/bridge"
	"github.com/galleguillosdavid-coder/Ip7-C/core/overlay"
	"github.com/galleguillosdavid-coder/Ip7-C/core/p2p"
	"github.com/galleguillosdavid-coder/Ip7-C/core/protocol"
)

// Version global de IPv7-IEU
const Version = "2.2.5"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error fatal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// --- Flags de configuración ---
	role       := flag.String("role", "master", "Rol del nodo: master o node")
	remoteIP   := flag.String("remote", "", "IP del nodo remoto")
	port       := flag.Int("port", 7777, "Puerto UDP del túnel")
	remotePort := flag.Int("remote-port", 7777, "Puerto UDP del nodo remoto")
	ifaceName  := flag.String("iface", "ieu0", "Nombre del adaptador virtual")
	useTUN     := flag.Bool("tun", true, "Habilitar adaptador virtual TUN")
	targetDID  := flag.String("did", "", "Búsqueda P2P de DID (ej. did:ipv7:100)")
	kernelMode := flag.String("kernel", "embed", "Motor de Kernel OS: legacy | embed | kamikaze")
	apiPort    := flag.Int("api-port", 7780, "Puerto de la REST API (0 = desactivada)")
	mqttBroker := flag.String("mqtt", "", "URL del broker MQTT (vacío = desactivado)")
	updateVerify  := flag.String("update-verify", "sha256", "Método de verificación de actualización: sha256 | none")
	subPort       := flag.Uint("sub-port", 0, "Sub-puerto lógico de este nodo (0-65535). Ej: --sub-port 8080")
	deviceType    := flag.String("device", "unknown", "Tipo de dispositivo: router|nat|ap|server|nas|edge|desktop|notebook|mobile|wearable|smarttv|console|printer|camera|iot-sensor|actuator|smarthome|industrial|leo|starlink|geo|meo|vehicle|drone")
	bootstrapNode := flag.String("bootstrap", "", "Endpoint de arranque Web3 para Malla P2P (ej. 192.168.0.1:4000)")
	flag.Parse()

	// --- 0.1 Asistente Interactivo Plug & Play ---
	if len(os.Args) == 1 {
		fmt.Println("\n===================================================")
		fmt.Println("   IPv7-IEU - Asistente de Despliegue Rápido")
		fmt.Println("===================================================")
		fmt.Print("👉 Deseas operar como Nodo Primario [MASTER]? (M)\n👉 Deseas conectarte como Satélite  [NODO]?   (N)\n> ")
		var resp string
		fmt.Scanln(&resp)
		if resp == "n" || resp == "N" {
			*role = "node"
		} else {
			*role = "master"
		}
	}

	// --- 0. Inyecciones preliminares (OS-specific) ---
	ensureSecurityExceptions()
	extractWintun(*kernelMode)
	go CheckUpdate(Version, *updateVerify)

	fmt.Printf("🚀 [IPv7-IEU] Protocolo de Ignición v%s - Modo: %s\n", Version, *role)

	// --- Identidad del nodo ---
	var localNode *protocol.Node
	var virtualIP string
	if *role == "master" {
		addr := protocol.NewIPv7WithSubPort(56, 1, 100, uint16(*subPort))
		localNode = &protocol.Node{
			Name:    "Bmax-B4-Master",
			Address: addr,
			Latency: 35.0,
		}
		virtualIP = "10.7.7.1"
	} else {
		fmt.Println("🔍 Entrando en la Fase Web3: Sincronización DHT Autónoma...")
		
		// rand.Seed eliminado: deprecated desde Go 1.20 (el PRNG global se seedea automáticamente)
		satID := float64(100 + rand.Intn(900))
		ipSuffix := 2 + rand.Intn(250)
		
		if *remoteIP == "" && *bootstrapNode == "" {
			fmt.Print("Ingresa la IP del MASTER/Bootstrap (Destino): ")
			fmt.Scanln(remoteIP)
		}
		
		fmt.Printf("🎯 Asignación Pseudo-DHCP Autónoma: ID=%.0f, VIP=10.7.7.%d\n", satID, ipSuffix)

		localNode = &protocol.Node{
			Name:    fmt.Sprintf("Satelite-Dinamico-%.0f", satID),
			Address: protocol.NewIPv7WithSubPort(56, 1, satID, uint16(*subPort)),
			Latency: 2.0,
		}
		virtualIP = fmt.Sprintf("10.7.7.%d", ipSuffix)
	}

	// --- 1. Inicializar Túnel Overlay (con NAT traversal automático) ---
	tunnel, err := overlay.NewTunnel(localNode, *port, *remoteIP, *remotePort)
	if err != nil {
		return fmt.Errorf("fallo al inicializar túnel: %v", err)
	}
	fmt.Printf("📡 Identidad Local: %.2f | Puerto UDP: %d\n", localNode.Address.ResolvedIP, *port)

	// Aplicar perfil de dispositivo al túnel
	devProfile := protocol.GetDeviceProfile(protocol.ParseDeviceClass(*deviceType))
	tunnel.SetDeviceProfile(devProfile)

	// Registrar IP pública via STUN (en background)
	go func() {
		publicEP, err := overlay.DiscoverPublicEndpoint()
		if err != nil {
			fmt.Printf("⚠️ [STUN] No se pudo descubrir IP pública: %v\n", err)
			return
		}
		fmt.Printf("🌍 [STUN] Endpoint público descubierto: %s\n", publicEP)
		tunnel.SetPublicEndpoint(publicEP)
	}()

	// --- 2. Red Web3 (DHT P2P) con Task Budget adaptado al dispositivo ---
	dhtPort := *port + 1000
	myDID := fmt.Sprintf("did:ipv7:%.0f", localNode.Address.ResolvedIP)
	microDHT, err := p2p.NewMicroDHT(myDID, dhtPort)
	if err != nil {
		fmt.Printf("⚠️ No se pudo inicializar la malla DHT Web3: %v\n", err)
	} else {
		// Aplicar el presupuesto recomendado por el perfil de dispositivo
		microDHT.SetBudget(devProfile.RecommendedDHTBudget, devProfile.BudgetWindow)
		fmt.Printf("🌐 [Web3] DID: %s | DHT Puerto: %d | Budget: %d RPCs/%s\n",
			myDID, dhtPort, devProfile.RecommendedDHTBudget, devProfile.BudgetWindow)

		if *bootstrapNode != "" {
			fmt.Printf("🎯 Conectando a Bootstrap Node Kademlia Web3: %s\n", *bootstrapNode)
			microDHT.AddBootstrap(*bootstrapNode)
		} else if *remoteIP != "" {
			microDHT.AddBootstrap(fmt.Sprintf("%s:%d", *remoteIP, dhtPort))
		}

		// Anuncio periódico del endpoint en la telaraña real P2P
		go func() {
			for {
				publicValue := fmt.Sprintf("IP_PÚBLICA:%d", *port)
				if ep := tunnel.GetPublicEndpoint(); ep != "" {
					publicValue = ep
				}
				microDHT.Announce(myDID, publicValue)
				time.Sleep(2 * time.Minute)
			}
		}()

		// Resolución DID si se especificó target
		if *targetDID != "" {
			fmt.Printf("🔍 Buscando identidad %s en la telaraña Kademlia...\n", *targetDID)
			if resolvedIP := microDHT.Resolve(*targetDID); resolvedIP != "" {
				fmt.Printf("✅ DID Resuelto desde P2P! Nodo encontrado en -> %s\n", resolvedIP)
			} else {
				fmt.Printf("❌ Nodo satelital DID fuera de rango estocástico.\n")
			}
		}
	}

	// --- 3. Balizas PQC OOB (Out-Of-Band) ---
	go func() {
		for {
			time.Sleep(30 * time.Second)
			if tunnel.RemoteAddr != nil {
				beaconData := []byte("QUANTUM_HANDSHAKE_SYNC")
				sig := protocol.GenerateSignature(beaconData)
				qAddr := protocol.NewIPv7(1, 1, 0xFFFFFFFF)
				pkt := append(qAddr.SerializeHeader(), beaconData...)
				pkt = append(pkt, sig...)
				tunnel.Conn.WriteToUDP(pkt, tunnel.RemoteAddr)
			}
		}
	}()

	// --- 4. Adaptador Virtual TUN ---
	var vIf *adapter.IEUInterface
	if *useTUN {
		fmt.Printf("🛠️ Creando interfaz virtual %s (%s)... ", *ifaceName, virtualIP)
		vIf, err = adapter.NewIEUInterface(*ifaceName, virtualIP, "255.255.255.0")
		if err != nil {
			fmt.Printf("\n⚠️ ADVERTENCIA: No se pudo crear interfaz virtual: %v\n", err)
			fmt.Println("   El sistema seguirá operando en modo Overlay Application puro.")
		} else {
			fmt.Println("¡Éxito!")
			defer vIf.Close()

			// 🔓 Perforación Automática de Firewall de Windows
			if runtime.GOOS == "windows" {
				go func() {
					// Regla genérica para la subred y otra específica para ICMP (ping)
					exec.Command("netsh", "advfirewall", "firewall", "add", "rule", "name=IPv7-IEU", "dir=in", "action=allow", "protocol=ANY", "remoteip=10.7.7.0/24").Run()
					exec.Command("netsh", "advfirewall", "firewall", "add", "rule", "name=IPv7-IEU-ICMP", "dir=in", "action=allow", "protocol=icmpv4").Run()
				}()
			}

			// Bridge: OS -> IEU (Micro-slicing activo)
			go func() {
				for {
					packet, err := vIf.Read()
					if err != nil {
						fmt.Printf("❌ Error leyendo de %s: %v\n", *ifaceName, err)
						continue
					}
					if *remoteIP != "" {
						if len(packet) < 128 {
							tunnel.SendPriority(packet)
						} else {
							tunnel.SendStandard(packet)
						}
					}
				}
			}()
		}
	}

	// --- 5. Listener General (Bridge: IEU -> OS) ---
	go tunnel.Listen(func(remoteAddr protocol.IPv7Address, data []byte) {
		if vIf != nil {
			if err := vIf.Write(data); err != nil {
				fmt.Printf("❌ Error escribiendo en %s: %v\n", *ifaceName, err)
			}
		} else {
			fmt.Printf("\n📩 [IEU-Packet] Recibido desde ID: %.2f | Size: %d bytes\n",
				remoteAddr.ResolvedIP, len(data))
		}
	})

	// --- 6. Bridges de Interoperabilidad ---
	nodeInfo := &bridge.NodeInfo{
		DID:     myDID,
		Version: Version,
		Role:    *role,
		Node:    localNode,
		Tunnel:  tunnel,
		DHT:     microDHT,
		APIPort: *apiPort, // Para WoT descriptor dinámico
	}
	
	if *role == "master" {
		bridge.StartMasterEgress(tunnel)
		// Nuevo: Enrutador Global TCP (Puerto fijo 9777 para fácil enrutamiento P2P)
		egressTcpPort := 9777
		go bridge.StartQuantumEgressServer(fmt.Sprintf("0.0.0.0:%d", egressTcpPort))
	} else {
		bridge.StartSatelliteEgress(tunnel)
		// Nuevo: Gateway Proxy Local
		masterIp := *remoteIP
		if masterIp == "" {
			masterIp = "127.0.0.1" // Fallback local
		}
		masterTcpAddr := fmt.Sprintf("%s:9777", masterIp)
		go bridge.StartSocks5Server("0.0.0.0:1080", masterTcpAddr)
	}

	if *apiPort > 0 {
		go bridge.StartRESTAPI(nodeInfo, *apiPort)
		fmt.Printf("🔌 [REST API] Escuchando en http://127.0.0.1:%d/v1/\n", *apiPort)
	}

	if *mqttBroker != "" {
		go bridge.StartMQTTBridge(nodeInfo, *mqttBroker)
		fmt.Printf("📨 [MQTT] Bridge conectado a %s\n", *mqttBroker)
	}

	go bridge.StartCoAPProxy(nodeInfo, 5683)
	fmt.Println("📡 [CoAP] Proxy UDP activo en :5683")

	// -- 7. MicroDHT Operación Contínua P2P --
	// Descubrimiento nativo mantenido vía threads asintomáticos,
	// desvinculación total de bases de datos de recolección en nube como Firebase.
	fmt.Println("🔗 MicroDHT Kademlia administrando latencias y roles en aislamiento Web3 total.")

	// --- Cierre limpio ---
	fmt.Println("\n---------------------------------------------------------")
	fmt.Println("  Presiona CTRL+C para cerrar el túnel y liberar la red.")
	fmt.Println("---------------------------------------------------------")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n🛑 Cerrando protocolo IPv7-IEU...")
	if runtime.GOOS == "windows" {
		exec.Command("netsh", "advfirewall", "firewall", "delete", "rule", "name=IPv7-IEU").Run()
		exec.Command("netsh", "advfirewall", "firewall", "delete", "rule", "name=IPv7-IEU-ICMP").Run()
	}
	fmt.Printf("✅ Interfaz %s desconectada. Red restaurada.\n", *ifaceName)
	return nil
}
