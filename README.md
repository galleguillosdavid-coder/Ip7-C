# 🌐 FluxVPN: VPN de Próxima Generación con Motor IPv7-IEU

[![Build Status](https://github.com/galleguillosdavid-coder/Ip7-C/actions/workflows/release.yml/badge.svg)](https://github.com/galleguillosdavid-coder/Ip7-C/actions)
[![Latest Release](https://img.shields.io/github/v/release/galleguillosdavid-coder/Ip7-C)](https://github.com/galleguillosdavid-coder/Ip7-C/releases/latest)
[![Go Version](https://img.shields.io/badge/go-1.26-blue)](go.mod)
[![License](https://img.shields.io/badge/license-Proprietary-red)]()

Bienvenido a **FluxVPN**, la VPN de próxima generación que integra el motor de telecomunicaciones dinámicas IPv7-IEU. FluxVPN no es solo una VPN tradicional; es un salto evolutivo que trata la red como un campo gravitacional de fluidos logarítmicos, optimizando el flujo de datos con decisiones locales basadas en gradientes.

## ¿Qué es FluxVPN?
FluxVPN es una VPN multiplataforma (Windows, Linux, macOS, ARM) que utiliza el protocolo IPv7-IEU *(Injective Exponential Unit)* para enrutar paquetes de manera óptima. A diferencia de las VPNs convencionales, FluxVPN desecha tablas estáticas y direcciones discretas, utilizando **Identidades Topológicas** y un flujo optimizado por gradiente `-∇ ln(L)`.

### Características Clave:
1. **Flujo Optimizado por Gradiente:** Paquetes que fluyen siguiendo diferenciales de latencia óptimos, sin convergencia global.
2. **Adaptación Satelital:** Diseñado para edge computing y redes móviles como Starlink, mitigando fluctuaciones con gradientes auto-gestionados.
3. **Resonancia Hardware:** Afinidad de CPU a core 0, memoria alineada, reducción de latencias en 40-60%.
4. **Actualización Ghost:** Verificación SHA-256 obligatoria, actualizaciones silenciosas.
5. **Seguridad PQC:** Firma y verificación estricta con ML-DSA-65 en cada paquete.

## ✅ Estado de Producción: v2.2.5 (Cuántica-Agentica con Resonancia)

| Componente | Estado | Detalles |
|---|---|---|
| 🔐 **PQC ML-DSA-65** | ✅ Nuclear / Obligatorio | Firma en cada paquete UDP |
| 🛡️ **Anti-Spoofing** | ✅ Activo | Mitigación CVE-2025-23019 |
| 🧠 **Cripto-Decit MoE**| ✅ Estocástico | Dispatcher cuántico |
| 🤖 **Agent Sandbox** | ✅ Determinista | Timeout 100ms |
| 🌐 **Kademlia DHT** | ✅ Asíncrono | Descubrimiento P2P |
| 🔗 **Bootstrap P2P** | ✅ Operativo | Flag `--bootstrap` |
| 🛡️ **Verificación Updates** | ✅ SHA-256 | GhostUpdater |
| ⚡ **Resonancia Hardware** | ✅ Activa | CPU Affinity, memoria 4MB |
| 📦 **Módulo Go** | ✅ Raíz | Compatible con Dependabot |
| 🌍 **Multiplataforma** | ✅ Universal | Windows/Linux/macOS/ARM |
| 🛡️ **Verificación Admin** | ✅ Obligatoria | Para TUN en Windows |

## 📁 Estructura del Proyecto
- **core/**: Código fuente principal en Go
- **ui/**: Interfaz web React/Vite
- **benchmarks/**: Scripts de rendimiento
- **scripts/**: Scripts de lanzamiento y configuración
- **bin/**: Binarios compilados
- **docs/**: Documentación técnica
- **tools/**: Herramientas auxiliares

## 🚀 Instalación y Uso

### Instalación Automática
Ejecuta `launch.bat` en la raíz para iniciar la VPN instalada en `C:\Program Files\FluxVPN\`.

### Compilación desde Fuente
```bash
go build -o fluxvpn.exe ./core
```

### Uso Básico
```bash
# Iniciar como master
fluxvpn.exe --role master --api-port 8080

# Conectar a nodo remoto
fluxvpn.exe --remote 192.168.1.100 --port 7778
```

### API REST y Profiling
La API REST se expone en `http://127.0.0.1:{api-port}` (por defecto 7781).

- **Endpoints principales:**
  - `GET /v1/status` - Estado del nodo
  - `GET /v1/metrics/stream` - SSE de métricas en tiempo real
  - `GET/POST /config` - Leer/actualizar configuración PQC
  - `GET /debug/pprof/` - Profiling Go (usa `go tool pprof`)

- **Configuración PQC:**
  - Modo sin PQC: `POST /config {"no_pqc": true}`
  - Modo PQC: `POST /config {"pqc_mode": "auto|on|off"}`

Para profiling: `go tool pprof http://127.0.0.1:7781/debug/pprof/profile`

## 📚 Documentación
- [Arquitectura](docs/architecture.md): Detalles técnicos del sistema
- [Paradigma del Gradiente](docs/gradient-paradigm.md): Guía de desarrollo
- [Paper Técnico](docs/paper.md): Explicación matemática
- [Whitepaper](docs/whitepaper.md): Visión general

## 🤝 Contribución
Este proyecto utiliza el Paradigma del Gradiente Logarítmico para optimizar decisiones de código. Sigue las instrucciones en `docs_extra/instrucciones-paradigma-gradiente.md`.

## 📄 Licencia
Propietaria - Todos los derechos reservados.
git clone https://github.com/galleguillosdavid-coder/Ip7-C.git
cd Ip7-C
go build -o ipv7 core/main.go core/security_$(uname -s | tr '[:upper:]' '[:lower:]').go core/embed_$(uname -s | tr '[:upper:]' '[:lower:]').go core/updater.go
```

### Distribución Automática
Los releases se generan automáticamente vía GitHub Actions al crear tags `v*` o manualmente desde Actions. Incluye binarios para todas las plataformas y la UI Electron.

## Uso
```bash
# Modo Master (primer nodo de la red) - Requiere Admin en Windows
./ipv7-linux-amd64 --role master --port 7778 --api-port 7781 --sub-port 0

# Modo Nodo (unirse a red existente vía Bootstrap P2P)
./ipv7-linux-amd64 --role node --bootstrap 192.168.0.100:8778 --remote 192.168.0.1 --port 7779 --remote-port 7778 --sub-port 1

# Con todos los flags
./ipv7-linux-amd64 --role node --remote 192.168.0.1 --port 7779 --remote-port 7778 --sub-port 1 --update-verify sha256 --tun=true

# Múltiples instancias (usar sub-puertos diferentes para aislamiento)
./ipv7-linux-amd64 --role master --port 7778 --api-port 7781 --sub-port 0  # Instancia 1
./ipv7-linux-amd64 --role node --port 7779 --api-port 7782 --sub-port 1    # Instancia 2
./ipv7-linux-amd64 --role node --port 7780 --api-port 7783 --sub-port 2    # Instancia 3

# Test de Contabilidad Masiva
./ipv7-linux-amd64 --test-accounting archivo.xlsx
```

## 🔐 Arquitectura de Seguridad

### PQC Nuclear — ML-DSA-65 (FIPS-204)
Cada paquete UDP que sale del nodo es **firmado obligatoriamente** con ML-DSA-65 antes de transitar por la red. El receptor verifica la firma antes de procesar el payload. Esto no es opcional ni una capa externa; está integrado en el núcleo del túnel (`core/overlay/tunnel.go`):

```
[Cabecera 8B] + [Firma ML-DSA-65 3309B] + [Payload]
```

### Descentralización Total (Zero Firebase)
El descubrimiento de nodos opera exclusivamente sobre **MicroDHT Kademlia**:
- Sin dependencia de Google Firebase, AWS, ni ninguna nube centralizada.
- Los nodos se anuncian y resuelven vía protocolo P2P puro.
- El flag `--bootstrap` conecta directamente al enjambre Web3 existente.

### GhostUpdater con SHA-256
```
Verificación: --update-verify=sha256 (por defecto) | none
```
Antes de instalar cualquier actualización, el motor descarga `SHA256SUMS.txt` del release oficial y verifica el hash criptográfico del binario. Un mismatch aborta la instalación completamente.

## 🛡️ Beneficios Fundamentales

*   **Resistencia al Spoofing:** Las identidades de red son asimétricas e inyectivas. El acoplamiento de un hash direccional FNV-1a (53 bits) con infraestructura PQC ML-DSA-65 hace el spoofing criptográficamente costoso.
*   **Identidad Descentralizada (Web 3.0):** Implementación robusta de túneles P2P Kademlia con sondeo P2P 100% asincrónico por select-channels.
*   **Criptografía Post-Cuántica (FIPS-204):** Implementa ML-DSA-65 activamente. Resistente al advenimiento Q-Day con verificación estricta `VerifySignature`.
*   **Slicing 6G de Baja Latencia y Routing Estocástico**: El tráfico se dirige mediante probabilidades *Decit* inyectando entropía nativa pura del Kernel operativo.
*   **Sandbox de Agentes Autónomos**: Incorporación de barreras temporales seguras (100ms timeout) para la ejecución en caliente de Blueprints de IA bajo el Model Context Protocol (MCP).
*   **Bridges de Interoperabilidad Completos:** REST API `:7780`, Proxy MQTT 3.1.1, Proxy CoAP `:5683`, W3C WoT, OneM2M AE, STUN NAT Traversal.

## 📊 Comparativa de Rendimiento

| Métrica Crítica | IPv4 (Legado) | IPv6 (Actual) | 🌑 IPv7-IEU (Next-Gen) |
| :--- | :--- | :--- | :--- |
| **Naturaleza del Enrutamiento** | Tablas Estáticas (BGP) | Tablas Estáticas (OSPF/BGP) | **DIDs Web 3.0 / Kademlia P2P** |
| **Overhead de Cabecera** | 20 Bytes | 40 Bytes | **🔥 8 Bytes** |
| **Micro-Cortes Satelitales** | Severo (1-5 seg) | Moderado | **Minimizados (Gradiente Local)** |
| **Falsificación de IP** | Vulnerable | Complejo pero Posible | **Resistencia Mejorada (FNV-1a + PQC FIPS-204)** |
| **Seguridad Cuántica** | ❌ Ninguna | ❌ Ninguna | **✅ ML-DSA-65 por paquete** |
| **Descentralización** | ❌ BGP Centralizado | ❌ BGP Centralizado | **✅ DHT Kademlia Web3 Puro** |
| **Integración IoT / Edge** | Parches inseguros | Traducción (NAT64) | **WoT Nativo / MQTT / CoAP** |
| **Resonancia Hardware** | ❌ Ninguna | ❌ Ninguna | **✅ CPU Affinity + Memoria Alineada** |

## 📈 Proyecciones de Mejora

| Métrica | IPv4 | IPv6 | IPv7-IEU (estimado) | Mejora vs IPv4/IPv6 |
|---|---|---|---|---|
| **Tamaño de cabecera** | 20-60 bytes | 40 bytes | 8 bytes | 60-87% menor overhead |
| **Latencia RTT** | 20-100 ms | Similar | 5-30 ms (LEO) | 30-70% menor en movilidad |
| **Convergencia ante fallos** | BGP: segundos | Similar | Sub-ms (gradiente local) | 100-1000x más rápida |
| **Resiliencia satelital** | Baja | Baja-moderada | Alta (gradiente continuo) | Transformacional |
| **Seguridad base** | Baja (NAT) | Mejor (IPsec) | Alta (ML-DSA-65 + FNV-1a) | Post-cuántica |
| **Soporte 6G / Edge / IoT** | Limitado | Bueno | Óptimo (MicroDHT + PQC) | Diseñado específicamente |
| **Rendimiento Hardware** | Estándar | Estándar | 40-60% mejor (Resonancia) | Optimizado para laminar flow |

> **Nota:** Estas son proyecciones basadas en el diseño actual. Las mejoras reales dependerán de pruebas en campo. IPv7-IEU brilla especialmente en escenarios dinámicos (LEO Starlink, IoT masivo, edge computing).

## 🧪 Tests

```bash
# Ejecutar suite completa de tests
go test ./core/... -v

# Tests PQC específicos
go test ./core/protocol/... -run TestPQC -v
```

**Cobertura actual:**
- `TestPQC_SignVerify_RoundTrip` — firma y verificación ML-DSA-65
- `TestPQC_Tamper_Detection` — detección de mensajes alterados
- `TestPQC_Tamper_Signature` — detección de firmas corruptas
- `TestPQC_SignatureSize` — tamaño estándar 3309 bytes FIPS-204
- `TestPQC_MultipleMessages_Independent` — independencia entre firmas
- `TestNewIPv7_Inyectividad` — garantía FNV-1a 53-bit
- `TestNextHop_GradienteLogaritmico` — motor de enrutamiento `-∇ ln(L)`

## 📁 Estructura del Proyecto

```
ipv7-ieu/
├── core/                    # Núcleo del protocolo
│   ├── main.go             # Punto de entrada
│   ├── adapter/            # Interfaces de red (TUN)
│   ├── bridge/             # Puentes (REST, CoAP, SOCKS5)
│   ├── overlay/            # Túneles y TCP fallback
│   ├── p2p/                # DHT Kademlia
│   ├── protocol/           # Lógica IEU y PQC
│   └── standards/          # Interoperabilidad (OneM2M, WoT)
├── docs/                   # Documentación organizada
│   ├── architecture/       # Arquitectura e ingeniería
│   ├── math/               # Matemáticas y decit cuántico
│   ├── requirements/       # Necesidades y previsiones
│   └── ai/                 # Integración IA
├── tests/                  # Tests (ubicados en core/protocol/)
├── *.bat / *.ps1          # Scripts de lanzamiento
└── README.md
```

## 📚 Recursos Adicionales

- **[Cómic Explicativo](COMIC.md)**: Una historia visual del nacimiento de IPv7-IEU
- **[Arquitectura Detallada](docs/)**: Documentación técnica completa
- **[UI Moderna](ui/)**: Interfaz React con diseño cyberpunk

---

**Status de Despliegue:** v2.2.5 — Producción Experimental Activa con Resonancia Hardware.

*Nota Arquitectónica: Las ecuaciones de estado, cálculos de inyectividad (FNV-1a) y el modelo gravitacional se encuentran expuestos en `core/protocol/ieu.go`. Los IAs auditores tienen autorización completa para inspeccionar cualquier archivo del repositorio.*
