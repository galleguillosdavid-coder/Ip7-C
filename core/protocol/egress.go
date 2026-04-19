package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Constantes para Egress (SubPort 443)
const (
	EgressSubPort = 443
	MaxChunkSize  = 1200 // Para caber en un UDP MTU seguro (1500)
)

// Tipos de mensaje Egress
const (
	EgressOpReq   = 0x01
	EgressOpChunk = 0x02
	EgressOpDone  = 0x03
)

// EgressRequest es el formato del paquete que el Satélite envía al Master
func BuildEgressRequest(url string) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(EgressOpReq)
	urlBytes := []byte(url)
	binary.Write(buf, binary.BigEndian, uint16(len(urlBytes)))
	buf.Write(urlBytes)
	return buf.Bytes()
}

// ParseEgressRequest extrae la URL del payload
func ParseEgressRequest(data []byte) (string, error) {
	if len(data) < 3 || data[0] != EgressOpReq {
		return "", fmt.Errorf("paquete EgressReq inválido")
	}
	urlLen := binary.BigEndian.Uint16(data[1:3])
	if len(data) < int(3+urlLen) {
		return "", fmt.Errorf("tamaño EgressReq corrupto")
	}
	return string(data[3 : 3+urlLen]), nil
}

// BuildEgressChunk crea un fragmento de datos del archivo
func BuildEgressChunk(seq uint32, chunkData []byte) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(EgressOpChunk)
	binary.Write(buf, binary.BigEndian, seq)
	buf.Write(chunkData)
	return buf.Bytes()
}

// BuildEgressDone señaliza el fin de la transmisión
func BuildEgressDone() []byte {
	return []byte{EgressOpDone}
}
