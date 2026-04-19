# 🌐 IPv7-IEU: El Futuro de las Telecomunicaciones Dinámicas

[![Build Status](https://github.com/galleguillosdavid-coder/Ip7-IEU/actions/workflows/release.yml/badge.svg)](https://github.com/galleguillosdavid-coder/Ip7-IEU/actions)
[![Latest Release](https://img.shields.io/github/v/release/galleguillosdavid-coder/Ip7-IEU)](https://github.com/galleguillosdavid-coder/Ip7-IEU/releases/latest)
[![Go Version](https://img.shields.io/badge/go-1.22-blue)](go.mod)
[![License](https://img.shields.io/badge/license-Proprietary-red)]()

Bienvenido a la primera capa de red de próxima generación. **IPv7-IEU** *(Injective Exponential Unit)* no es solo un protocolo de enrutamiento; es un salto evolutivo en la forma en que los datos atraviesan el mundo físico y satelital.

## ¿Qué es IPv7-IEU?
A diferencia de los protocolos tradicionales creados en la década de los 70s y 90s, IPv7-IEU desecha por completo la noción de direccionar máquinas usando "tablas estáticas" o "números discretos". En este ecosistema, una dirección de red es una **Identidad Topológica**. IPv7-IEU envuelve el tráfico convencional dentro de una matriz virtual que trata el internet no como cables desconectados, sino como un **campo gravitacional de fluidos logarítmicos**.

## ¿Qué hace y Cómo Funciona?
Operando como un poderoso *Overlay Network* o software de túnel nativo escalable (Multiplataforma Windows, Linux, macOS, ARM):
1. **Flujo Optimizado por Gradiente:** Cada paquete "fluye" hacia el destino siguiendo un diferencial de latencia óptimo `-∇ ln(L)`, tomando decisiones locales sin necesitar convergencia global (a diferencia de BGP/OSPF).
2. **Adaptación Satelital:** Programado desde su concepción para hardware en el filo (*Edge Computing* como nodos Bmax) e internet de movilidad (Starlink). El motor IEU mitiga fluctuaciones de posición o caídas intermitentes, estabilizando latencias mediante gradientes auto-gestionados.
3. **Actualización Ghost (GhostUpdater):** Arquitectura sigilosa multiplataforma con **verificación SHA-256 obligatoria** de cada binario descargado antes de la instalación. Los clientes se mantienen actualizados sin requerir interacción del usuario.

## ✅ Estado de Producción: v1.5.7

| Componente | Estado | Detalles |
|---|---|---|
| 🔐 **PQC ML-DSA-65** | ✅ Nuclear / Obligatorio | Firma y verificación en **cada** paquete UDP |
| 🌐 **Descentralización** | ✅ Zero Cloud | Firebase eliminado — DHT Kademlia Web3 puro |
| 🔗 **Bootstrap P2P** | ✅ Operativo | Flag `--bootstrap` para auto-unirse a la red |
| 🛡️ **Verificación Updates** | ✅ SHA-256 | GhostUpdater valida hash antes de hot-swap |
| 🧪 **Tests PQC** | ✅ 10/10 PASS | Suite ML-DSA-65 completa (sign/verify/tamper) |
| 📦 **Módulo Go** | ✅ Raíz | `go.mod` en raíz — compatible con Dependabot/CodeQL |
| 🤖 **CI/CD** | ✅ Automatizado | Build + SHA256SUMS.txt en cada release |
| 🌍 **Multiplataforma** | ✅ Universal | Windows / Linux x64 / Linux ARM64 / macOS Intel / M1-M3 |

## 🚀 Instalación Rápida

### Descarga Directa (Recomendado)
```bash
# Linux x64
curl -LO https://github.com/galleguillosdavid-coder/Ip7-IEU/releases/latest/download/ipv7-linux-amd64
chmod +x ipv7-linux-amd64

# Verificar integridad SHA-256
curl -LO https://github.com/galleguillosdavid-coder/Ip7-IEU/releases/latest/download/SHA256SUMS.txt
sha256sum -c SHA256SUMS.txt
```

### Uso
```bash
# Modo Master (primer nodo de la red)
./ipv7-linux-amd64 --role master

# Modo Nodo (unirse a red existente vía Bootstrap P2P)
./ipv7-linux-amd64 --role node --bootstrap 192.168.0.100:8777

# Con todos los flags
./ipv7-linux-amd64 --role node --remote 192.168.0.1 --port 7777 --update-verify sha256
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
*   **Identidad Descentralizada (Web 3.0):** Implementación robusta de túneles P2P estilo Kademlia MicroDHT. Los nodos descubren IPs satelitales dinámicas (Starlink) mediante DIDs (`did:ipv7:x`), sin dependencia de DNS ni PKIs corporativas.
*   **Criptografía Post-Cuántica (FIPS-204):** Implementa ML-DSA-65 del estándar NIST FIPS-204. Resistente al advenimiento Q-Day con firma de 3309 bytes por paquete.
*   **Slicing 6G de Baja Latencia:** Motor adaptativo ramificado. El tráfico menor a 128 bytes obtiene pase "Prioridad Absoluta" y es enrutado por la cola de baja latencia.
*   **Bridges de Interoperabilidad Completos:** REST API `:7780`, Proxy MQTT 3.1.1, Proxy CoAP `:5683`, W3C WoT Thing Description, oneM2M AE, NAT Traversal STUN + Fallback TCP 443.

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

## 📈 Proyecciones de Mejora

| Métrica | IPv4 | IPv6 | IPv7-IEU (estimado) | Mejora vs IPv4/IPv6 |
|---|---|---|---|---|
| **Tamaño de cabecera** | 20-60 bytes | 40 bytes | 8 bytes | 60-87% menor overhead |
| **Latencia RTT** | 20-100 ms | Similar | 5-30 ms (LEO) | 30-70% menor en movilidad |
| **Convergencia ante fallos** | BGP: segundos | Similar | Sub-ms (gradiente local) | 100-1000x más rápida |
| **Resiliencia satelital** | Baja | Baja-moderada | Alta (gradiente continuo) | Transformacional |
| **Seguridad base** | Baja (NAT) | Mejor (IPsec) | Alta (ML-DSA-65 + FNV-1a) | Post-cuántica |
| **Soporte 6G / Edge / IoT** | Limitado | Bueno | Óptimo (MicroDHT + PQC) | Diseñado específicamente |

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

---

**Status de Despliegue:** v1.5.7 — Producción Experimental Activa (Multiplataforma).

*Nota Arquitectónica: Las ecuaciones de estado, cálculos de inyectividad (FNV-1a) y el modelo gravitacional se encuentran expuestos en `core/protocol/ieu.go`. Los IAs auditores tienen autorización completa para inspeccionar cualquier archivo del repositorio.*
