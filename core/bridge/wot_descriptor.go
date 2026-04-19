package bridge

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// WoTProperty describe una propiedad observable del nodo en formato W3C WoT
type WoTProperty struct {
	Type        string `json:"type"`
	Unit        string `json:"unit,omitempty"`
	ReadOnly    bool   `json:"readOnly"`
	Observable  bool   `json:"observable"`
	Description string `json:"description"`
}

// WoTAction describe una acción invocable en el nodo
type WoTAction struct {
	Description string      `json:"description"`
	Input       interface{} `json:"input,omitempty"`
	Output      interface{} `json:"output,omitempty"`
}

// WoTForm describe cómo acceder a una propiedad o acción
type WoTForm struct {
	Href        string `json:"href"`
	ContentType string `json:"contentType"`
	Op          string `json:"op"`
}

// ThingDescription es el documento JSON-LD W3C WoT Thing Description del nodo IEU.
// Permite que plataformas como Home Assistant, AWS IoT Core, FIWARE, Eclipse Thingweb
// descubran y consuman este nodo sin configuración manual.
type ThingDescription struct {
	Context    []interface{}             `json:"@context"`
	ID         string                    `json:"id"`
	Title      string                    `json:"title"`
	Type       []string                  `json:"@type"`
	Version    map[string]string         `json:"version"`
	Created    string                    `json:"created"`
	Modified   string                    `json:"modified"`
	Support    string                    `json:"support"`
	Security   []string                  `json:"security"`
	SecurityDef map[string]interface{}   `json:"securityDefinitions"`
	Properties map[string]WoTPropFull    `json:"properties"`
	Actions    map[string]WoTActionFull  `json:"actions"`
}

// WoTPropFull incluye el esquema completo de la propiedad con sus formas de acceso
type WoTPropFull struct {
	Type        string    `json:"type"`
	Unit        string    `json:"unit,omitempty"`
	ReadOnly    bool      `json:"readOnly"`
	Observable  bool      `json:"observable"`
	Description string    `json:"description"`
	Forms       []WoTForm `json:"forms"`
}

// WoTActionFull incluye el esquema completo de la acción con input/output
type WoTActionFull struct {
	Description string      `json:"description"`
	Input       interface{} `json:"input,omitempty"`
	Output      interface{} `json:"output,omitempty"`
	Forms       []WoTForm   `json:"forms"`
}

// GenerateThingDescription escribe el Thing Description JSON-LD al writer dado.
// Usar con Content-Type: application/td+json (RFC del W3C WoT)
func GenerateThingDescription(w io.Writer, info *NodeInfo) {
	now := time.Now().UTC().Format(time.RFC3339)
	port := info.APIPort
	if port == 0 {
		port = 7780 // fallback al puerto por defecto
	}
	baseURL := fmt.Sprintf("http://127.0.0.1:%d", port)

	td := ThingDescription{
		Context: []interface{}{
			"https://www.w3.org/2019/wot/td/v1",
			map[string]string{
				"ipv7": "https://ipv7-ieu.protocol/vocab#",
			},
		},
		ID:    info.DID,
		Title: fmt.Sprintf("IPv7-IEU Node [%s]", info.Role),
		Type:  []string{"Thing", "ipv7:NetworkNode"},
		Version: map[string]string{
			"instance": info.Version,
		},
		Created:  now,
		Modified: now,
		Support:  "https://github.com/galleguillosdavid-coder/Ip7-IEU-Releases",
		Security: []string{"nosec_sc"},
		SecurityDef: map[string]interface{}{
			"nosec_sc": map[string]string{"scheme": "nosec"},
		},
		Properties: map[string]WoTPropFull{
			"latency": {
				Type:        "number",
				Unit:        "millisecond",
				ReadOnly:    true,
				Observable:  true,
				Description: "Latencia actual del nodo en milisegundos",
				Forms: []WoTForm{
					{Href: baseURL + "/v1/status", ContentType: "application/json", Op: "readproperty"},
				},
			},
			"did": {
				Type:        "string",
				ReadOnly:    true,
				Observable:  false,
				Description: "Identificador Descentralizado W3C del nodo IPv7-IEU",
				Forms: []WoTForm{
					{Href: baseURL + "/v1/status", ContentType: "application/json", Op: "readproperty"},
				},
			},
			"peers": {
				Type:        "array",
				ReadOnly:    true,
				Observable:  true,
				Description: "Lista de peers conocidos en la malla Kademlia DHT",
				Forms: []WoTForm{
					{Href: baseURL + "/v1/peers", ContentType: "application/json", Op: "readproperty"},
				},
			},
			"publicEndpoint": {
				Type:        "string",
				ReadOnly:    true,
				Observable:  false,
				Description: "Endpoint público IP:Puerto descubierto via STUN",
				Forms: []WoTForm{
					{Href: baseURL + "/v1/status", ContentType: "application/json", Op: "readproperty"},
				},
			},
			"pqcPubKey": {
				Type:        "string",
				ReadOnly:    true,
				Observable:  false,
				Description: "Clave pública ML-DSA-65 (FIPS 204) del nodo para verificación PQC",
				Forms: []WoTForm{
					{Href: baseURL + "/v1/pqc/pubkey", ContentType: "application/json", Op: "readproperty"},
				},
			},
		},
		Actions: map[string]WoTActionFull{
			"send": {
				Description: "Envía un payload a un nodo destino identificado por DID",
				Input: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"did":      map[string]string{"type": "string"},
						"payload":  map[string]string{"type": "string"},
						"priority": map[string]string{"type": "boolean"},
					},
					"required": []string{"did", "payload"},
				},
				Output: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"status": map[string]string{"type": "string"},
					},
				},
				Forms: []WoTForm{
					{Href: baseURL + "/v1/send", ContentType: "application/json", Op: "invokeaction"},
				},
			},
		},
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(td); err != nil {
		fmt.Fprintf(w, `{"error":"no se pudo generar Thing Description: %v"}`, err)
	}
}
