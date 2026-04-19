package bridge

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
)

// MQTTBridge conecta un broker MQTT con la red IPv7-IEU.
// Traduce mensajes MQTT a paquetes IEU y viceversa sin dependencias externas
// de cliente MQTT (implementación wire-level mínima del protocolo MQTT 3.1.1).

// StartMQTTBridge conecta a un broker MQTT y establece el bridge bidireccional.
// brokerURL formato: "tcp://host:1883" o "host:1883"
func StartMQTTBridge(info *NodeInfo, brokerURL string) {
	// Normalizar URL
	addr := strings.TrimPrefix(brokerURL, "tcp://")
	addr = strings.TrimPrefix(addr, "mqtt://")

	for {
		if err := runMQTTBridge(info, addr); err != nil {
			fmt.Printf("⚠️ [MQTT] Conexión perdida: %v — Reconectando en 10s...\n", err)
			time.Sleep(10 * time.Second)
		}
	}
}

func runMQTTBridge(info *NodeInfo, addr string) error {
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return fmt.Errorf("no se pudo conectar al broker %s: %v", addr, err)
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	clientID := strings.ReplaceAll(info.DID, ":", "_")

	// --- MQTT CONNECT packet (protocolo 3.1.1) ---
	if err := mqttConnect(conn, clientID); err != nil {
		return fmt.Errorf("error en CONNECT: %v", err)
	}

	// Leer CONNACK
	connack := make([]byte, 4)
	if _, err := conn.Read(connack); err != nil {
		return fmt.Errorf("error leyendo CONNACK: %v", err)
	}
	if connack[0] != 0x20 || connack[3] != 0x00 {
		return fmt.Errorf("broker rechazó la conexión, código: %d", connack[3])
	}
	fmt.Printf("✅ [MQTT] Conectado al broker %s como %s\n", addr, clientID)

	// --- SUBSCRIBE al topic de entrada ---
	inboundTopic := fmt.Sprintf("ipv7/%s/send", info.DID)
	if err := mqttSubscribe(conn, inboundTopic, 1); err != nil {
		return fmt.Errorf("error en SUBSCRIBE: %v", err)
	}

	// Reset deadline para operación continua
	conn.SetDeadline(time.Time{})

	// --- Loop de lectura de mensajes del broker ---
	for {
		conn.SetDeadline(time.Now().Add(5 * time.Minute))
		packetType, payload, err := mqttReadPacket(conn)
		if err != nil {
			return fmt.Errorf("error leyendo paquete MQTT: %v", err)
		}

		switch packetType {
		case 0x30: // PUBLISH
			topicLen := int(payload[0])<<8 | int(payload[1])
			topic := string(payload[2 : 2+topicLen])
			message := payload[2+topicLen:]

			fmt.Printf("📨 [MQTT->IEU] Topic: %s | %d bytes\n", topic, len(message))

			// Encapsular mensaje MQTT como paquete IEU
			if len(message) < 128 {
				info.Tunnel.SendPriority(message)
			} else {
				info.Tunnel.SendStandard(message)
			}

		case 0xD0: // PINGRESP — mantener keepalive
			// No-op
		}
	}
}

// PublishToMQTT publica un mensaje IEU recibido al topic de salida del nodo.
// Se llama desde el handler del tunnel cuando llega un paquete destinado a este nodo.
func PublishToMQTT(conn net.Conn, info *NodeInfo, data []byte) {
	outboundTopic := fmt.Sprintf("ipv7/%s/inbox", info.DID)
	mqttPublish(conn, outboundTopic, data, 0)
}

// --- Helpers de bajo nivel para el protocolo MQTT 3.1.1 ---

func mqttConnect(conn net.Conn, clientID string) error {
	idBytes := []byte(clientID)
	// Variable header: Protocol Name "MQTT", Level 4, Flags 0x02 (Clean Session), Keepalive 60s
	varHeader := []byte{
		0x00, 0x04, 'M', 'Q', 'T', 'T', // Protocol Name
		0x04,       // Protocol Level (3.1.1)
		0x02,       // Connect Flags: Clean Session
		0x00, 0x3C, // Keep Alive: 60 seconds
	}
	payload := append([]byte{0x00, byte(len(idBytes))}, idBytes...)
	remaining := append(varHeader, payload...)

	packet := []byte{0x10} // CONNECT type
	packet = append(packet, mqttEncodeLength(len(remaining))...)
	packet = append(packet, remaining...)

	_, err := conn.Write(packet)
	return err
}

func mqttSubscribe(conn net.Conn, topic string, qos byte) error {
	topicBytes := []byte(topic)
	// Packet identifier = 0x0001
	payload := []byte{0x00, 0x01}
	payload = append(payload, 0x00, byte(len(topicBytes)))
	payload = append(payload, topicBytes...)
	payload = append(payload, qos)

	packet := []byte{0x82} // SUBSCRIBE type + reserved bits
	packet = append(packet, mqttEncodeLength(len(payload))...)
	packet = append(packet, payload...)

	_, err := conn.Write(packet)
	return err
}

func mqttPublish(conn net.Conn, topic string, data []byte, qos byte) error {
	topicBytes := []byte(topic)
	varHeader := append([]byte{0x00, byte(len(topicBytes))}, topicBytes...)
	payload := append(varHeader, data...)

	packet := []byte{0x30} // PUBLISH, QoS 0, no retain
	packet = append(packet, mqttEncodeLength(len(payload))...)
	packet = append(packet, payload...)

	_, err := conn.Write(packet)
	return err
}

func mqttReadPacket(conn net.Conn) (byte, []byte, error) {
	header := make([]byte, 1)
	if _, err := conn.Read(header); err != nil {
		return 0, nil, err
	}

	length, err := mqttDecodeLength(conn)
	if err != nil {
		return 0, nil, err
	}

	body := make([]byte, length)
	if length > 0 {
		if _, err := conn.Read(body); err != nil {
			return 0, nil, err
		}
	}
	return header[0], body, nil
}

func mqttEncodeLength(length int) []byte {
	var encoded []byte
	for {
		digit := byte(length % 128)
		length /= 128
		if length > 0 {
			digit |= 0x80
		}
		encoded = append(encoded, digit)
		if length == 0 {
			break
		}
	}
	return encoded
}

func mqttDecodeLength(conn net.Conn) (int, error) {
	multiplier := 1
	value := 0
	buf := make([]byte, 1)
	for {
		if _, err := conn.Read(buf); err != nil {
			return 0, err
		}
		value += int(buf[0]&0x7F) * multiplier
		multiplier *= 128
		if buf[0]&0x80 == 0 {
			break
		}
		if multiplier > 128*128*128 {
			return 0, fmt.Errorf("campo de longitud MQTT inválido")
		}
	}
	return value, nil
}

// GetPeerList expone la lista de peers del MicroDHT para la API REST
func (info *NodeInfo) GetPeerListJSON() []byte {
	if info.DHT == nil {
		return []byte("[]")
	}
	peers := info.DHT.GetPeerList()
	b, _ := json.Marshal(peers)
	return b
}
