# 🌐 IPv7-IEU — Guía de Interoperabilidad

Este documento explica cómo conectar sistemas externos a la red IPv7-IEU usando los bridges de compatibilidad incluidos en `v1.4.0-Universal`.

---

## Instalación por Plataforma

| Plataforma | Binario | Requisito |
|---|---|---|
| Windows x64 | `ipv7-windows-amd64.exe` | Correr como Administrador |
| Linux x64 | `ipv7-linux-amd64` | `sudo` o `CAP_NET_ADMIN` |
| Linux ARM64 | `ipv7-linux-arm64` | Raspberry Pi 3/4/5 |
| macOS Intel | `ipv7-darwin-amd64` | `sudo` |
| macOS Apple Silicon | `ipv7-darwin-arm64` | `sudo` |

```bash
# Linux — dar permiso de ejecución y lanzar
chmod +x ipv7-linux-amd64
sudo ./ipv7-linux-amd64 -role=master -port=7777

# macOS
sudo ./ipv7-darwin-arm64 -role=master
```

---

## Flags de Configuración

```
-role         master | node         (default: master)
-remote       <IP del peer>         IP del nodo remoto
-port         <puerto>              Puerto UDP principal (default: 7777)
-remote-port  <puerto>              Puerto UDP del peer (default: 7777)
-iface        <nombre>              Nombre del adaptador virtual (default: ieu0)
-tun          true | false          Habilitar adaptador TUN (default: true)
-did          <did:ipv7:XXX>        DID a buscar en la red P2P
-kernel       embed|legacy|kamikaze Motor de kernel (solo Windows)
-api-port     <puerto>              Puerto REST API (default: 7780, 0=desactivar)
-mqtt         <url>                 URL broker MQTT (vacío=desactivar)
```

---

## REST API — `http://127.0.0.1:7780/v1/`

La REST API permite que cualquier aplicación (web, móvil, scripts) controle e inspeccione el nodo.

### Endpoints

#### `GET /v1/status`
Estado completo del nodo.

```bash
curl http://127.0.0.1:7780/v1/status
```

```json
{
  "did": "did:ipv7:5600",
  "version": "1.4.0-Universal",
  "role": "master",
  "resolved_ip": 5600,
  "latency_ms": 35,
  "public_endpoint": "203.0.113.42:7777",
  "timestamp": "2026-04-17T20:00:00Z"
}
```

#### `GET /v1/peers`
Lista de peers en la malla DHT P2P.

```bash
curl http://127.0.0.1:7780/v1/peers
```

#### `POST /v1/send`
Enviar un mensaje a otro nodo identificado por DID.

```bash
curl -X POST http://127.0.0.1:7780/v1/send \
  -H "Content-Type: application/json" \
  -d '{"did":"did:ipv7:101","payload":"hola desde REST","priority":true}'
```

#### `GET /v1/pqc/pubkey`
Clave pública ML-DSA-65 del nodo.

#### `GET /v1/wot`
Thing Description W3C WoT en formato JSON-LD.

#### `GET /v1/health`
Health check simple para load balancers y monitoreo.

---

## MQTT Bridge

El bridge MQTT connect to any standard MQTT 3.1.1 broker (Mosquitto, HiveMQ, EMQX, AWS IoT Core, etc.).

### Activar

```bash
./ipv7-linux-amd64 -mqtt tcp://localhost:1883
```

### Topics

| Topic | Dirección | Descripción |
|---|---|---|
| `ipv7/{mi_did}/send` | → IEU | Publicar aquí para enviar al nodo IEU |
| `ipv7/{mi_did}/inbox` | IEU → | Mensajes recibidos por el nodo se publican aquí |

### Ejemplo con Mosquitto

```bash
# En una terminal: lanzar broker local
mosquitto -v

# En otra terminal: publicar un mensaje al nodo IEU
mosquitto_pub -t "ipv7/did_ipv7_5600/send" -m "hola sensor"

# Suscribirse a los mensajes que llegan al nodo
mosquitto_sub -t "ipv7/did_ipv7_5600/inbox"
```

---

## CoAP Bridge (IoT dispositivos a batería)

El proxy CoAP escucha en UDP puerto `5683` (estándar RFC 7252).

### Recursos disponibles

| URI-Path | Método | Descripción |
|---|---|---|
| `/ieu/status` | GET | Estado del nodo en formato texto compacto |
| `/ieu/did` | GET | DID del nodo |
| `/ieu/send` | POST | Enviar payload al nodo IEU |

### Ejemplo con coap-client

```bash
# Instalar: apt install libcoap3-bin
coap-get coap://localhost:5683/ieu/status
coap-post coap://localhost:5683/ieu/send -e "sensor_data=42"
```

### Ejemplo con Python (aiocoap)

```python
import asyncio
import aiocoap

async def main():
    protocol = await aiocoap.Context.create_client_context()
    request = aiocoap.Message(
        code=aiocoap.POST,
        uri='coap://127.0.0.1/ieu/send',
        payload=b'temperatura=25.3'
    )
    response = await protocol.request(request).response
    print(f'Respuesta: {response.payload}')

asyncio.run(main())
```

---

## W3C Web of Things (WoT)

Cada nodo auto-genera su Thing Description (TD) en formato JSON-LD.

```bash
curl http://127.0.0.1:7780/v1/wot | python -m json.tool
```

### Integración con Eclipse Thingweb

```bash
# Descubrir el nodo IEU desde Thingweb Node-WoT
node wot-server.js --td-url http://127.0.0.1:7780/v1/wot
```

### Integración con Home Assistant

Agregar en `configuration.yaml`:
```yaml
rest:
  - resource: http://127.0.0.1:7780/v1/status
    scan_interval: 30
    sensor:
      - name: "IPv7 Latencia"
        value_template: "{{ value_json.latency_ms }}"
        unit_of_measurement: "ms"
```

---

## NAT Traversal

IPv7-IEU v1.4.0+ descubre automáticamente tu IP pública via STUN y realiza hole-punching para conectar nodos detrás de NAT sin configuración manual.

### Flujo automático

1. Al iniciar, el nodo consulta `stun.l.google.com:19302`
2. Descubre su `IP_pública:Puerto` real
3. Anuncia este endpoint en la red DHT Kademlia
4. Cuando un peer quiere conectar, ambos lados intercambian endpoints via DHT y realizan hole-punch simultáneo

### Fallback TCP

Si UDP está bloqueado por el firewall, el sistema automáticamente intenta:
1. TCP en el mismo puerto configurado
2. TCP en puerto 8443
3. TCP en puerto 443 (camouflage como HTTPS)

---

## oneM2M — Gestión Industrial

El nodo se registra automáticamente como Application Entity (AE) en plataformas oneM2M.

Para usar un CSE interno corporativo, modificar `core/standards/onem2m.go`:

```go
const oneM2MEndpoint = "http://tu-cse-interno.empresa.com:8282/~/mn-cse/mn-name"
```

---

## Arquitectura de Compatibilidad

```
                     ┌─────────────────────────────┐
                     │       IPv7-IEU Node          │
                     │                              │
  REST Clients ──────┤  REST API :7780              │
  Web / Mobile       │                              │
                     │  MQTT Bridge ────────────────┼──→ Broker MQTT
  IoT Sensors ───────┤  CoAP Proxy :5683            │    (Mosquitto/HiveMQ)
  (batería)          │                              │
                     │  W3C WoT /v1/wot             │
  Home Assistant ────┤  Thing Description           │
  AWS IoT / FIWARE   │                              │
                     │  oneM2M AE Registration      │
  Eclipse OM2M ──────┤  (industrial fleet mgmt)     │
                     │                              │
                     │  UDP Tunnel (IEU) :7777       │
  Otros nodos ───────┤  STUN + NAT Traversal        │
  IPv7-IEU           │  TCP Fallback :443           │
                     │                              │
                     │  Kademlia DHT :8777          │
  P2P Network ───────┤  did:ipv7: resolution        │
                     └─────────────────────────────┘
```
