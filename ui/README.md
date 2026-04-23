# FluxVPN IPv7 - UI

Interfaz de usuario para FluxVPN IPv7, una VPN P2P con criptografía post-cuántica.

## Características

- **Interfaz Moderna**: UI React con diseño cyberpunk
- **Telemetría en Tiempo Real**: SSE para métricas live
- **Dashboard de Gráficos**: Visualización de latencia y throughput
- **Modo Experto**: Detalles técnicos avanzados
- **App Nativa**: Empaquetada con Electron

## Desarrollo

```bash
npm install
npm run dev          # Vite dev server
npm run electron-dev # Electron con hot reload
npm run dist         # Build para distribución
```

## Distribución

Los releases se generan automáticamente vía GitHub Actions al crear tags `v*`:

```bash
git tag v1.0.0
git push origin v1.0.0
```

O manualmente desde la pestaña Actions.

Los binarios incluyen:
- `ipv7` (daemon Go) para Linux, Windows, macOS
- `FluxVPN IPv7` (app Electron) para cada plataforma

## Arquitectura

- **Frontend**: React + Vite
- **Backend**: Go daemon (IPv7-IEU)
- **Empaquetado**: Electron + electron-builder
- **Criptografía**: ML-DSA-65 PQC + HMAC SHA-256

## Protocolo

IPv7-IEU v1.5.7 con enrutamiento gradiente logarítmico (-∇ ln(L))
