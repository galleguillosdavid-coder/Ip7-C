package bridge

import (
	"github.com/galleguillosdavid-coder/Ip7-C/core/overlay"
	"github.com/galleguillosdavid-coder/Ip7-C/core/p2p"
	"github.com/galleguillosdavid-coder/Ip7-C/core/protocol"
)

// NodeInfo es el contexto compartido entre todos los bridges.
// Permite a la REST API, MQTT bridge y CoAP proxy acceder al estado del nodo IEU.
type NodeInfo struct {
	DID     string
	Version string
	Role    string
	Node    *protocol.Node
	Tunnel  *overlay.Tunnel
	DHT     *p2p.MicroDHT
	APIPort int // Puerto de la REST API para WoT descriptor dinámico
}
