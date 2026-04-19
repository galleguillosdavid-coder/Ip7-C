package standards

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// oneM2MEndpoint es el endpoint base de la plataforma oneM2M pública de prueba.
// En producción, reemplazar con el CSE-Base propio o un broker corporativo.
const oneM2MEndpoint = "http://onem2m.mcloud.edu.gr:8282/~/mn-cse/mn-name"

// oneM2MAERequest representa la estructura de registro de un Application Entity (AE) en oneM2M
type oneM2MAERequest struct {
	M2MAE AEBody `json:"m2m:ae"`
}

// AEBody contiene los campos del Application Entity según TS-0001 oneM2M
type AEBody struct {
	ResourceName   string   `json:"rn"`          // Resource Name (identificador)
	AppName        string   `json:"apn"`         // Application Name
	AppID          string   `json:"api"`         // Application ID
	RequestReachability bool `json:"rr"`         // Si el AE es alcanzable (true)
	SupportedReleaseVersions []string `json:"srv"` // Versiones oneM2M soportadas
	Labels         []string `json:"lbl"`         // Etiquetas para búsqueda semántica
}

// RegisterAsOneM2MAE registra el nodo IPv7-IEU como un Application Entity (AE)
// en una plataforma oneM2M, permitiendo que sistemas de gestión de dispositivos
// industriales gestionen el nodo remotamente.
func RegisterAsOneM2MAE(did string, version string) {
	// Esperar a que el sistema esté completamente inicializado
	time.Sleep(5 * time.Second)

	fmt.Println("🏭 [oneM2M] Intentando registro como Application Entity...")

	resourceName := sanitizeDID(did)

	body := oneM2MAERequest{
		M2MAE: AEBody{
			ResourceName:    resourceName,
			AppName:         "IPv7-IEU",
			AppID:           "Napp7IEU001",
			RequestReachability: true,
			SupportedReleaseVersions: []string{"3", "2a"},
			Labels: []string{
				"protocol:ipv7",
				"version:" + version,
				"pqc:ml-dsa-65",
				"network:p2p-kademlia",
				"did:" + did,
			},
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("⚠️ [oneM2M] Error serializando AE body: %v\n", err)
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", oneM2MEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("⚠️ [oneM2M] Error creando request: %v\n", err)
		return
	}

	// Headers oneM2M obligatorios (TS-0009)
	req.Header.Set("Content-Type", "application/json;ty=2") // ty=2 = AE resource type
	req.Header.Set("X-M2M-Origin", resourceName)
	req.Header.Set("X-M2M-RI", fmt.Sprintf("ipv7-reg-%d", time.Now().UnixMilli()))
	req.Header.Set("X-M2M-RVI", "3") // Release Version Indicator

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("⚠️ [oneM2M] No se pudo conectar al CSE: %v (continuando sin registro)\n", err)
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 201:
		fmt.Printf("✅ [oneM2M] Nodo %s registrado exitosamente como AE\n", resourceName)
	case 409:
		fmt.Printf("ℹ️ [oneM2M] AE %s ya existe en el CSE (registro previo)\n", resourceName)
	default:
		fmt.Printf("⚠️ [oneM2M] Respuesta inesperada del CSE: HTTP %d\n", resp.StatusCode)
	}
}

// sanitizeDID convierte un DID como "did:ipv7:5600" a un resource name válido "did_ipv7_5600"
func sanitizeDID(did string) string {
	result := ""
	for _, c := range did {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' {
			result += string(c)
		} else {
			result += "_"
		}
	}
	return result
}
