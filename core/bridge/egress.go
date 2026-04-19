package bridge

import (
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/galleguillosdavid-coder/Ip7-C/core/overlay"
	"github.com/galleguillosdavid-coder/Ip7-C/core/protocol"
)

// egressHTTPClient con timeout para no colgarse indefinidamente en descargas
var egressHTTPClient = &http.Client{Timeout: 60 * time.Second}

// StartMasterEgress registra un listener en el Master en el SubPort 443
// para atender requests de descarga remotos desde Satélites.
func StartMasterEgress(tunnel *overlay.Tunnel) {
	tunnel.RegisterSubPort(protocol.EgressSubPort, func(remoteAddr protocol.IPv7Address, data []byte) {
		url, err := protocol.ParseEgressRequest(data)
		if err != nil {
			fmt.Println("❌ [Egress Master] Error parseando request:", err)
			return
		}

		fmt.Printf("🌐 [Egress Master] Satellite %.0f pide descargar: %s\n", remoteAddr.ResolvedIP, url)

		go func(targetUrl string, cltAddr protocol.IPv7Address) {
			resp, err := egressHTTPClient.Get(targetUrl)
			if err != nil {
				fmt.Println("❌ [Egress Master] Error HTTP GET:", err)
				tunnel.SendSubPort(cltAddr, protocol.EgressSubPort, []byte("ERROR: "+err.Error()))
				return
			}
			defer resp.Body.Close()

			seq := uint32(0)
			buf := make([]byte, protocol.MaxChunkSize)

			for {
				n, err := resp.Body.Read(buf)
				if n > 0 {
					chunk := protocol.BuildEgressChunk(seq, buf[:n])
					tunnel.SendSubPort(cltAddr, protocol.EgressSubPort, chunk)
					GlobalTelemetry.BytesSent.Add(int64(len(chunk)))
					GlobalTelemetry.PacketsSent.Add(1)
					seq++
				}
				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Println("❌ [Egress Master] Error leyendo body:", err)
					break
				}
			}
			// Enviar fin de transmisión
			tunnel.SendSubPort(cltAddr, protocol.EgressSubPort, protocol.BuildEgressDone())
			fmt.Printf("✅ [Egress Master] Archivo enviado al Satellite %.0f (%d chunks)\n", cltAddr.ResolvedIP, seq)
		}(url, remoteAddr)
	})
}

// ── Satellite Egress con buffer de reordenamiento ────────────────────────────────────────────

// satelliteDownload administra el estado de una descarga activa en el satélite.
// Reordena chunks out-of-order usando el número de secuencia del protocolo Egress.
// Esto resuelve el bug donde los chunks UDP podían llegar fuera de orden y corrompían el archivo.
type satelliteDownload struct {
	mu      sync.Mutex
	file    *os.File
	pending map[uint32][]byte // chunks fuera de orden en espera
	nextSeq uint32            // siguiente seq esperado para escritura ordenada
}

func newSatelliteDownload(f *os.File) *satelliteDownload {
	return &satelliteDownload{
		file:    f,
		pending: make(map[uint32][]byte),
	}
}

// writeChunk almacena el chunk con su seq y vacía los que ya están en orden.
func (dl *satelliteDownload) writeChunk(seq uint32, data []byte) {
	dl.mu.Lock()
	defer dl.mu.Unlock()
	dl.pending[seq] = append([]byte(nil), data...)
	// Vaciar todos los chunks consecutivos ya disponibles
	for {
		chunk, ok := dl.pending[dl.nextSeq]
		if !ok {
			break
		}
		if dl.file != nil {
			dl.file.Write(chunk)
		}
		delete(dl.pending, dl.nextSeq)
		dl.nextSeq++
	}
}

func (dl *satelliteDownload) close() {
	dl.mu.Lock()
	defer dl.mu.Unlock()
	if dl.file != nil {
		dl.file.Close()
		dl.file = nil
	}
}

// dlMu protege el puntero global activeDownload contra acceso concurrente
var dlMu sync.Mutex
var activeDownload *satelliteDownload

// StartSatelliteEgress registra un listener en el Satélite para recibir los chunks
func StartSatelliteEgress(tunnel *overlay.Tunnel) {
	tunnel.RegisterSubPort(protocol.EgressSubPort, func(remoteAddr protocol.IPv7Address, data []byte) {
		if len(data) == 0 {
			return
		}

		switch data[0] {
		case protocol.EgressOpChunk:
			if len(data) > 5 {
				// Extraer número de secuencia para reordenamiento garantizado
				seq := binary.BigEndian.Uint32(data[1:5])
				chunkData := data[5:]
				GlobalTelemetry.BytesReceived.Add(int64(len(data)))
				GlobalTelemetry.PacketsReceived.Add(1)
				dlMu.Lock()
				dl := activeDownload
				dlMu.Unlock()
				if dl != nil {
					dl.writeChunk(seq, chunkData)
				}
			}
		case protocol.EgressOpDone:
			fmt.Println("✅ [Egress Satellite] Descarga Cósmica Completada via IPv7!")
			dlMu.Lock()
			dl := activeDownload
			activeDownload = nil
			dlMu.Unlock()
			if dl != nil {
				dl.close()
			}
		}
	})
}

// TriggerSatelliteDownload manda la orden al Master para iniciar una descarga proxy
func TriggerSatelliteDownload(tunnel *overlay.Tunnel, masterDID string, url string) error {
	fmt.Printf("🚀 [Egress Satellite] Solicitando descarga autónoma: %s\n", url)

	// Crear archivo local temporal con os.CreateTemp (ioutil.TempFile deprecated desde Go 1.16)
	f, err := os.CreateTemp(".", "ipv7_ai_model_*.bin")
	if err != nil {
		return err
	}
	fmt.Printf("📦 [Egress Satellite] Guardando en: %s\n", f.Name())

	dlMu.Lock()
	if activeDownload != nil {
		activeDownload.close() // Cancelar descarga anterior si existía
	}
	activeDownload = newSatelliteDownload(f)
	dlMu.Unlock()

	req := protocol.BuildEgressRequest(url)

	// Rutear vía UDP al endpoint remoto configurado (topología 1-a-1)
	masterAddr := protocol.NewIPv7(56, 1, 100)
	err = tunnel.SendSubPort(masterAddr, protocol.EgressSubPort, req)
	return err
}
