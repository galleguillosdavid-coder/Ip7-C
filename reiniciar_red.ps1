Write-Host "🚧 Compilando la nueva versión IPv7-IEU (Cuántica-Agentica con Resonancia)..."
go build -o ipv7.exe ./core
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Error compilando. Revisa el código."
    exit
}

Write-Host "🛑 Desconectando la versión vieja de Internet (ipv7.exe) sin romper puertos..."
Stop-Process -Name ipv7 -Force -ErrorAction SilentlyContinue

Write-Host "⏳ Pausa crítica de 3 segundos para asegurar la liberación del puerto UDP y el TCP..."
Start-Sleep -Seconds 3

Write-Host "🌐 Reactivando Internet: Iniciando Master IPv7-IEU en su propia ventana..."
Start-Process -FilePath ".\ipv7.exe" -ArgumentList "--role master --port 7778 --api-port 7781 --tun=true" -Verb RunAs

Write-Host "⏳ Pausa de 2 segundos para inicializar DHT y Sandbox Agentico del Master..."
Start-Sleep -Seconds 2

Write-Host "🛰️ Conectando Satélite Cliente (Node) en su propia ventana..."
Start-Process -FilePath ".\ipv7.exe" -ArgumentList "--role node --port 7779 --api-port 7782 --tun=true --bootstrap 127.0.0.1:8778 --remote 127.0.0.1 --remote-port 7778" -Verb RunAs

Write-Host "✅ Red reestablecida con consolas interactivas completas y resonancia hardware."
