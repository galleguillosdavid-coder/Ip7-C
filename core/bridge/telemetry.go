package bridge

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/galleguillosdavid-coder/Ip7-IEU/core/protocol"
)

// ─── Telemetry Hub (Previsiones.md §Visibilidad End-to-End / interfaz.md §SSE) ──────────────
//
// "La degradación parcial e intermitente es algorítmicamente más difícil de diagnosticar
//  que un apagón total, exigiendo telemetría semántica." — Previsiones.md
//
// TelemetryHub acumula métricas en tiempo real del nodo y las sirve vía SSE a cualquier
// dashboard conectado. Los contadores son atómicos para ser seguros en goroutines.

// TelemetryHub es el registro centralizado de métricas del nodo.
type TelemetryHub struct {
	mu sync.RWMutex

	// Contadores absolutos (atómicos)
	PacketsReceived  atomic.Int64
	PacketsSent      atomic.Int64
	PacketsDropped   atomic.Int64
	RSErrorsDetected atomic.Int64
	BytesReceived    atomic.Int64
	BytesSent        atomic.Int64

	// Estado del MoE Dispatcher
	ActiveExpert    atomic.Int32 // 0=Latency 1=Bulk 2=Satellite
	DHTBudgetUsed   atomic.Int32
	DHTBudgetMax    atomic.Int32

	// Predicción de latencia Fokker-Planck
	LatencyPredMean   float64
	LatencyPredStdDev float64
	LatencyRisk       float64 // Probabilidad de degradación

	// Peers activos (copia del DHT)
	PeerAddrs []string

	// Contadores de ventana (para cálculo de tasa /seg)
	lastWindow time.Time
	prevRx     int64
	prevTx     int64
	RxPerSec   float64
	TxPerSec   float64

	// Canal para notificar a clientes SSE
	subscribers map[chan []byte]struct{}
	subMu       sync.Mutex
}

// GlobalTelemetry es el hub singleton accesible desde cualquier módulo.
var GlobalTelemetry = &TelemetryHub{
	subscribers: make(map[chan []byte]struct{}),
	lastWindow:  time.Now(),
}

// Subscribe registra un canal para recibir eventos SSE.
func (h *TelemetryHub) Subscribe() chan []byte {
	ch := make(chan []byte, 32)
	h.subMu.Lock()
	h.subscribers[ch] = struct{}{}
	h.subMu.Unlock()
	return ch
}

// Unsubscribe elimina y cierra el canal de un suscriptor SSE.
func (h *TelemetryHub) Unsubscribe(ch chan []byte) {
	h.subMu.Lock()
	delete(h.subscribers, ch)
	h.subMu.Unlock()
	close(ch)
}

// Publish serializa las métricas actuales y las envía a todos los suscriptores SSE.
func (h *TelemetryHub) Publish() {
	h.mu.Lock()
	now := time.Now()
	elapsed := now.Sub(h.lastWindow).Seconds()
	if elapsed > 0 {
		rxNow := h.PacketsReceived.Load()
		txNow := h.PacketsSent.Load()
		h.RxPerSec = float64(rxNow-h.prevRx) / elapsed
		h.TxPerSec = float64(txNow-h.prevTx) / elapsed
		h.prevRx = rxNow
		h.prevTx = txNow
		h.lastWindow = now
	}
	snapshot := map[string]interface{}{
		"ts":                 now.UTC().Format(time.RFC3339Nano),
		"packets_rx":         h.PacketsReceived.Load(),
		"packets_tx":         h.PacketsSent.Load(),
		"packets_dropped":    h.PacketsDropped.Load(),
		"rs_errors":          h.RSErrorsDetected.Load(),
		"bytes_rx":           h.BytesReceived.Load(),
		"bytes_tx":           h.BytesSent.Load(),
		"rx_per_sec":         h.RxPerSec,
		"tx_per_sec":         h.TxPerSec,
		"moe_expert":         expertName(int(h.ActiveExpert.Load())),
		"dht_budget_used":    h.DHTBudgetUsed.Load(),
		"dht_budget_max":     h.DHTBudgetMax.Load(),
		"latency_pred_mean":  h.LatencyPredMean,
		"latency_pred_std":   h.LatencyPredStdDev,
		"latency_risk":       h.LatencyRisk,
		"peers":              h.PeerAddrs,
		"peer_count":         len(h.PeerAddrs),
	}
	h.mu.Unlock()

	b, _ := json.Marshal(snapshot)
	msg := append([]byte("data: "), b...)
	msg = append(msg, '\n', '\n')

	h.subMu.Lock()
	for ch := range h.subscribers {
		select {
		case ch <- msg:
		default: // canal lleno, descartar (evita bloqueo)
		}
	}
	h.subMu.Unlock()
}

// StartPublishLoop publica métricas cada segundo en background.
func (h *TelemetryHub) StartPublishLoop() {
	go func() {
		for range time.NewTicker(time.Second).C {
			h.Publish()
		}
	}()
}

// SetLatencyPrediction actualiza la predicción de Fokker-Planck.
func (h *TelemetryHub) SetLatencyPrediction(mean, std, risk float64) {
	h.mu.Lock()
	h.LatencyPredMean = mean
	h.LatencyPredStdDev = std
	h.LatencyRisk = risk
	h.mu.Unlock()
}

// SetPeers actualiza la lista de peers de la DHT.
func (h *TelemetryHub) SetPeers(peers []string) {
	h.mu.Lock()
	h.PeerAddrs = make([]string, len(peers))
	copy(h.PeerAddrs, peers)
	h.mu.Unlock()
}

func expertName(e int) string {
	switch e {
	case 0:
		return "Latency"
	case 1:
		return "Bulk"
	case 2:
		return "Satellite"
	default:
		return "Unknown"
	}
}

// ─── SSE Handler: GET /v1/metrics/stream ─────────────────────────────────────────────────────
// Envía métricas en tiempo real como Server-Sent Events (SSE).
// Compatible con EventSource del navegador sin configuración adicional.
func handleMetricsSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE no soportado por este servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // Para nginx proxies

	ch := GlobalTelemetry.Subscribe()
	defer GlobalTelemetry.Unsubscribe(ch)

	// Enviar evento inicial inmediatamente
	fmt.Fprintf(w, "event: connected\ndata: {\"msg\":\"IPv7-IEU Telemetry Stream activo\"}\n\n")
	flusher.Flush()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			w.Write(msg)
			flusher.Flush()
		}
	}
}

// ─── Snapshot Handler: GET /v1/metrics ───────────────────────────────────────────────────────
// Devuelve un snapshot JSON puntual de las métricas (para polling si SSE no está disponible).
func handleMetricsSnapshot(w http.ResponseWriter, _ *http.Request) {
	h := GlobalTelemetry
	h.mu.RLock()
	data := map[string]interface{}{
		"ts":                time.Now().UTC().Format(time.RFC3339),
		"packets_rx":        h.PacketsReceived.Load(),
		"packets_tx":        h.PacketsSent.Load(),
		"packets_dropped":   h.PacketsDropped.Load(),
		"rs_errors":         h.RSErrorsDetected.Load(),
		"bytes_rx":          h.BytesReceived.Load(),
		"bytes_tx":          h.BytesSent.Load(),
		"rx_per_sec":        h.RxPerSec,
		"tx_per_sec":        h.TxPerSec,
		"moe_expert":        expertName(int(h.ActiveExpert.Load())),
		"dht_budget_used":   h.DHTBudgetUsed.Load(),
		"dht_budget_max":    h.DHTBudgetMax.Load(),
		"latency_pred_mean": h.LatencyPredMean,
		"latency_risk":      h.LatencyRisk,
		"peer_count":        len(h.PeerAddrs),
	}
	h.mu.RUnlock()
	jsonResponse(w, data, http.StatusOK)
}

// ─── Reset Handler: POST /v1/metrics/reset ───────────────────────────────────────────────────
// Pone a cero los contadores de la sesión (botón "Resetear Métricas" del dashboard).
func handleMetricsReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	h := GlobalTelemetry
	h.PacketsReceived.Store(0)
	h.PacketsSent.Store(0)
	h.PacketsDropped.Store(0)
	h.RSErrorsDetected.Store(0)
	h.BytesReceived.Store(0)
	h.BytesSent.Store(0)
	h.mu.Lock()
	h.prevRx = 0
	h.prevTx = 0
	h.lastWindow = time.Now()
	h.mu.Unlock()
	jsonResponse(w, map[string]string{"status": "contadores reiniciados"}, http.StatusOK)
}

// TrafficClassFromDevice retorna la TrafficClass recomendada para un DeviceClass.
func TrafficClassFromDevice(device protocol.DeviceClass) protocol.TrafficClass {
	switch device {
	case protocol.DeviceIndustrial, protocol.DeviceIoTActuator, protocol.DeviceDrone:
		return protocol.TC_REALTIME
	case protocol.DeviceIoTSensor, protocol.DeviceSmartHome, protocol.DevicePrinter:
		return protocol.TC_BACKGROUND
	case protocol.DeviceNAS, protocol.DeviceServer:
		return protocol.TC_BULK
	default:
		return protocol.TC_CONTROL
	}
}
