package bridge

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/galleguillosdavid-coder/Ip7-IEU/core/overlay"
	"github.com/galleguillosdavid-coder/Ip7-IEU/core/protocol"
)

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
			resp, err := http.Get(targetUrl)
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

var currentDownload *os.File

// StartSatelliteEgress registra un listener en el Satélite para recibir los chunks
func StartSatelliteEgress(tunnel *overlay.Tunnel) {
	tunnel.RegisterSubPort(protocol.EgressSubPort, func(remoteAddr protocol.IPv7Address, data []byte) {
		if len(data) == 0 {
			return
		}
		
		switch data[0] {
		case protocol.EgressOpChunk:
			if len(data) > 5 {
				chunkData := data[5:]
				if currentDownload != nil {
					currentDownload.Write(chunkData)
					// Publish telemetría simulada
					GlobalTelemetry.BytesReceived.Add(int64(len(data)))
					GlobalTelemetry.PacketsReceived.Add(1)
				}
			}
		case protocol.EgressOpDone:
			fmt.Println("✅ [Egress Satellite] Descarga Cósmica Completada via IPv7!")
			if currentDownload != nil {
				currentDownload.Close()
				currentDownload = nil
			}
		}
	})
}

// TriggerSatelliteDownload manda la orden al Master
func TriggerSatelliteDownload(tunnel *overlay.Tunnel, masterDID string, url string) error {
	fmt.Printf("🚀 [Egress Satellite] Solicitando descarga autónoma: %s\n", url)
	
	// Crear archivo local temporal
	var err error
	currentDownload, err = ioutil.TempFile(".", "ipv7_ai_model_*.bin")
	if err != nil {
		return err
	}
	fmt.Printf("📦 [Egress Satellite] Guardando en: %s\n", currentDownload.Name())

	req := protocol.BuildEgressRequest(url)
	
	// Como no tenemos el address del master directo sin resolver, lo forzamos con un fake address 
	// porque tunnel.Send lo rutea vía UDP remote endpoint de todas formas en 1 to 1.
	masterAddr := protocol.NewIPv7(56, 1, 100)
	err = tunnel.SendSubPort(masterAddr, protocol.EgressSubPort, req)
	if err != nil {
		return err
	}
	return nil
}
