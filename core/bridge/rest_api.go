package bridge

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/galleguillosdavid-coder/Ip7-C/core/protocol"
)

//go:embed dashboard.html
var dashboardFS embed.FS

// StartRESTAPI arranca la API REST + Dashboard en 127.0.0.1:{port}.
// Solo expone endpoints en loopback por seguridad.
func StartRESTAPI(info *NodeInfo, port int) {
	// Arrancar el loop de publicación de métricas SSE
	GlobalTelemetry.StartPublishLoop()

	mux := http.NewServeMux()

	// ── Dashboard ──────────────────────────────────────────────────────────
	// GET / → sirve el dashboard HTML desde el binario embebido
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "/dashboard" {
			http.NotFound(w, r)
			return
		}
		data, err := dashboardFS.ReadFile("dashboard.html")
		if err != nil {
			http.Error(w, "Dashboard no disponible", 500)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	})

	// ── API v1 ────────────────────────────────────────────────────────────
	mux.HandleFunc("/v1/status", func(w http.ResponseWriter, r *http.Request) {
		handleStatus(w, r, info)
	})
	mux.HandleFunc("/v1/peers", func(w http.ResponseWriter, r *http.Request) {
		handlePeers(w, r, info)
	})
	mux.HandleFunc("/v1/send", func(w http.ResponseWriter, r *http.Request) {
		handleSend(w, r, info)
	})
	mux.HandleFunc("/v1/pqc/pubkey", func(w http.ResponseWriter, r *http.Request) {
		handlePQCPubKey(w, r)
	})
	mux.HandleFunc("/v1/egress", func(w http.ResponseWriter, r *http.Request) {
		handleEgress(w, r, info)
	})
	mux.HandleFunc("/v1/wot", func(w http.ResponseWriter, r *http.Request) {
		handleWoTDescriptor(w, r, info)
	})
	mux.HandleFunc("/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// ── Métricas (Previsiones.md §Visibilidad End-to-End / interfaz.md) ───
	// GET  /v1/metrics         → snapshot JSON puntual
	// GET  /v1/metrics/stream  → SSE streaming en tiempo real
	// POST /v1/metrics/reset   → reinicia contadores de sesión
	mux.HandleFunc("/v1/metrics", func(w http.ResponseWriter, r *http.Request) {
		handleMetricsSnapshot(w, r)
	})
	mux.HandleFunc("/v1/metrics/stream", func(w http.ResponseWriter, r *http.Request) {
		handleMetricsSSE(w, r)
	})
	mux.HandleFunc("/v1/metrics/reset", func(w http.ResponseWriter, r *http.Request) {
		handleMetricsReset(w, r)
	})

	// ── QoS / Slices (Necesidades.md §QoS / Previsiones.md §Network Slicing) ─
	// GET /v1/slices → Catálogo de slice profiles para todos los dispositivos
	mux.HandleFunc("/v1/slices", func(w http.ResponseWriter, r *http.Request) {
		handleSlices(w, r, info)
	})

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	url := "http://" + addr
	fmt.Printf("🖥️  [Dashboard] Abre en navegador: %s\n", url)
	
	// Abrir automáticamente la URL del Dashboard en el navegador por defecto
	if runtime.GOOS == "windows" {
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	} else if runtime.GOOS == "darwin" {
		exec.Command("open", url).Start()
	} else {
		exec.Command("xdg-open", url).Start()
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      corsMiddleware(mux),
		ReadTimeout:  0, // SSE requiere sin timeout de read
		WriteTimeout: 0,
		IdleTimeout:  120 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("❌ [REST API] Error fatal: %v\n", err)
	}
}

// GET /v1/status — Estado completo del nodo
func handleStatus(w http.ResponseWriter, r *http.Request, info *NodeInfo) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	resp := map[string]interface{}{
		"did":             info.DID,
		"version":         info.Version,
		"role":            info.Role,
		"resolved_ip":     info.Node.Address.ResolvedIP,
		"latency_ms":      info.Node.Latency,
		"public_endpoint": info.Tunnel.GetPublicEndpoint(),
		"timestamp":       time.Now().UTC().Format(time.RFC3339),
	}
	jsonResponse(w, resp, http.StatusOK)
}

// GET /v1/peers — Lista de peers conocidos en la DHT
func handlePeers(w http.ResponseWriter, r *http.Request, info *NodeInfo) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	if info.DHT == nil {
		jsonResponse(w, map[string]interface{}{"peers": []string{}, "error": "DHT no inicializada"}, 200)
		return
	}
	peers := info.DHT.GetPeerList()

	// Actualizar telemetría con la lista de peers
	GlobalTelemetry.SetPeers(peers)

	jsonResponse(w, map[string]interface{}{
		"count": len(peers),
		"peers": peers,
	}, http.StatusOK)
}

// POST /v1/send — Enviar payload a un DID destino
// Body: {"did": "did:ipv7:101", "payload": "base64_or_text", "priority": true, "traffic_class": "realtime"}
func handleSend(w http.ResponseWriter, r *http.Request, info *NodeInfo) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	var body struct {
		DID          string `json:"did"`
		Payload      string `json:"payload"`
		Priority     bool   `json:"priority"`
		TrafficClass string `json:"traffic_class"` // "control","realtime","bulk","background"
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Bad Request: "+err.Error(), 400)
		return
	}
	if body.DID == "" || body.Payload == "" {
		http.Error(w, "Bad Request: did y payload son requeridos", 400)
		return
	}

	data := []byte(body.Payload)

	// Seleccionar SendPriorityOnSubPort según TrafficClass
	tc := protocol.TC_BULK
	switch body.TrafficClass {
	case "control":
		tc = protocol.TC_CONTROL
	case "realtime":
		tc = protocol.TC_REALTIME
	case "background":
		tc = protocol.TC_BACKGROUND
	}
	subPort := protocol.SubPortWithTC(tc, 0)

	if body.Priority || tc == protocol.TC_CONTROL || tc == protocol.TC_REALTIME {
		info.Tunnel.SendPriorityOnSubPort(data, subPort)
	} else {
		info.Tunnel.SendStandard(data)
	}

	// Actualizar contadores de telemetría
	GlobalTelemetry.PacketsSent.Add(1)
	GlobalTelemetry.BytesSent.Add(int64(len(data)))

	jsonResponse(w, map[string]string{
		"status":        "enviado",
		"did":           body.DID,
		"traffic_class": body.TrafficClass,
	}, http.StatusOK)
}

// GET /v1/pqc/pubkey — Clave pública ML-DSA-65 del nodo en hex
func handlePQCPubKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	pub := protocol.GetPublicKey()
	if pub == nil {
		http.Error(w, "Clave pública no disponible", 503)
		return
	}
	jsonResponse(w, map[string]interface{}{
		"algorithm":  "ML-DSA-65 (FIPS 204)",
		"pubkey_len": len(pub),
		"pubkey_hex": fmt.Sprintf("%x", pub[:32]) + "...", // preview primeros 32 bytes
	}, http.StatusOK)
}

// GET /v1/wot — Thing Description W3C WoT del nodo
func handleWoTDescriptor(w http.ResponseWriter, r *http.Request, info *NodeInfo) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "application/td+json")
	GenerateThingDescription(w, info)
}

// POST /v1/egress — Activa la descarga cuántica hacia Master
func handleEgress(w http.ResponseWriter, r *http.Request, info *NodeInfo) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Query parameter 'url' is required", 400)
		return
	}

	masterDID := r.URL.Query().Get("master_did")
	
	err := TriggerSatelliteDownload(info.Tunnel, masterDID, url)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	
	jsonResponse(w, map[string]interface{}{
		"status": "download_started",
		"target": url,
	}, 200)
}

// GET /v1/slices — Catálogo de Network Slice Profiles disponibles
func handleSlices(w http.ResponseWriter, r *http.Request, info *NodeInfo) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	type SliceSummary struct {
		Key         string  `json:"key"`
		Name        string  `json:"name"`
		Device      string  `json:"device"`
		TC          string  `json:"traffic_class"`
		MaxLatency  float64 `json:"max_latency_ms"`
		RSSymbols   int     `json:"rs_symbols"`
		DHTBudget   int     `json:"dht_budget"`
		DeltaHeader bool    `json:"delta_header"`
		Description string  `json:"description"`
	}
	tcName := func(tc protocol.TrafficClass) string {
		switch tc {
		case protocol.TC_CONTROL:
			return "control"
		case protocol.TC_REALTIME:
			return "realtime"
		case protocol.TC_BULK:
			return "bulk"
		default:
			return "background"
		}
	}
	result := make([]SliceSummary, 0, len(protocol.SliceProfiles))
	for key, sp := range protocol.SliceProfiles {
		dp := protocol.GetDeviceProfile(sp.Device)
		result = append(result, SliceSummary{
			Key:         key,
			Name:        sp.Name,
			Device:      dp.Name,
			TC:          tcName(sp.TC),
			MaxLatency:  sp.MaxLatencyMs,
			RSSymbols:   sp.RSSymbols,
			DHTBudget:   sp.DHTBudget,
			DeltaHeader: sp.DeltaHeader,
			Description: sp.Description,
		})
	}
	jsonResponse(w, map[string]interface{}{
		"count":  len(result),
		"slices": result,
	}, http.StatusOK)
}

// corsMiddleware permite peticiones CORS desde herramientas locales de desarrollo
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func jsonResponse(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
