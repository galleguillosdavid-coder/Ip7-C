Write-Host "🚧 Compilando la nueva versión IPv7-C (Cuántica-Agentica)..."
go build -o ipv7-c.exe ./core
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Error compilando. Revisa el código."
    exit
}

Write-Host "🛑 Desconectando la versión vieja de Internet (ipv7-node.exe) sin romper puertos..."
Stop-Process -Name ipv7-node -Force -ErrorAction SilentlyContinue
Stop-Process -Name ipv7-c -Force -ErrorAction SilentlyContinue

Write-Host "⏳ Pausa crítica de 3 segundos para asegurar la liberación del puerto UDP y el TCP..."
Start-Sleep -Seconds 3

Write-Host "🌐 Reactivando Internet: Iniciando Master IPv7-C en su propia ventana..."
Start-Process -FilePath ".\ipv7-c.exe" -ArgumentList "-role master -port 7777 -api-port 7780 -tun=false"

Write-Host "⏳ Pausa de 2 segundos para inicializar DHT y Sandbox Agentico del Master..."
Start-Sleep -Seconds 2

Write-Host "🛰️ Conectando Satélite Cliente (Node) en su propia ventana..."
Start-Process -FilePath ".\ipv7-c.exe" -ArgumentList "-role node -port 7778 -api-port 7781 -tun=false -bootstrap 127.0.0.1:8777 -remote 127.0.0.1"

Write-Host "✅ Red reestablecida con consolas interactivas completas."
