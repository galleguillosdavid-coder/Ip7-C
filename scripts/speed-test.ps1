# Script para prueba de velocidad real de IPv7-IEU
# Mide velocidad de descarga con y sin proxy SOCKS5

Write-Host "🚀 Prueba de Velocidad Real IPv7-IEU"
Write-Host "====================================="

# Archivo de prueba (1MB de datos aleatorios)
$url = "https://httpbin.org/bytes/1048576"  # 1MB

# Función para medir descarga
function Measure-Speed {
    param($useProxy, $label)
    
    Write-Host "`n📊 Midiendo: $label"
    
    if ($useProxy) {
        # Configurar proxy SOCKS5
        $env:HTTP_PROXY = "socks5://127.0.0.1:1080"
        $env:HTTPS_PROXY = "socks5://127.0.0.1:1080"
    } else {
        # Sin proxy
        Remove-Item Env:HTTP_PROXY -ErrorAction SilentlyContinue
        Remove-Item Env:HTTPS_PROXY -ErrorAction SilentlyContinue
    }
    
    try {
        $start = Get-Date
        Invoke-WebRequest -Uri $url -OutFile "test_speed.tmp" -TimeoutSec 30
        $end = Get-Date
        
        $duration = ($end - $start).TotalSeconds
        $fileSizeMB = (Get-Item "test_speed.tmp").Length / 1MB
        
        $speedMbps = ($fileSizeMB * 8) / $duration
        
        Write-Host "✅ Tamaño: $([math]::Round($fileSizeMB, 2)) MB"
        Write-Host "⏱️  Tiempo: $([math]::Round($duration, 2)) segundos"
        Write-Host "🚀 Velocidad: $([math]::Round($speedMbps, 2)) Mbps"
        
        return $speedMbps
    } catch {
        Write-Host "❌ Error: $($_.Exception.Message)"
        return 0
    } finally {
        Remove-Item "test_speed.tmp" -ErrorAction SilentlyContinue
    }
}

# Verificar que IPv7 esté corriendo
$ipv7Processes = Get-Process -Name "ipv7" -ErrorAction SilentlyContinue
if ($ipv7Processes.Count -lt 2) {
    Write-Host "⚠️  Advertencia: Se necesitan 2 procesos IPv7 corriendo (Master + Node)"
    Write-Host "Ejecuta Lanzar-Maestro.bat y Lanzar-Satelite.bat primero"
    exit 1
}

Write-Host "✅ IPv7 corriendo: $($ipv7Processes.Count) procesos"

# Prueba sin proxy (internet normal)
$speedNormal = Measure-Speed -useProxy $false -label "Sin IPv7 (Internet Normal)"

# Prueba con proxy IPv7
$speedIPv7 = Measure-Speed -useProxy $true -label "Con IPv7 (Proxy SOCKS5)"

# Resultados comparativos
Write-Host "`n📈 Resultados Comparativos"
Write-Host "=========================="
Write-Host "Sin IPv7: $([math]::Round($speedNormal, 2)) Mbps"
Write-Host "Con IPv7: $([math]::Round($speedIPv7, 2)) Mbps"

if ($speedIPv7 -gt 0 -and $speedNormal -gt 0) {
    $ratio = $speedIPv7 / $speedNormal
    $improvement = ($ratio - 1) * 100
    Write-Host "Mejora: $([math]::Round($improvement, 1))% $(if ($improvement -gt 0) {'más rápido'} else {'más lento'})"
    
    if ($ratio -gt 0.9) {
        Write-Host "🎉 Excelente: IPv7 mantiene velocidad comparable con resonancia hardware"
    } elseif ($ratio -gt 0.7) {
        Write-Host "👍 Bueno: IPv7 con buena performance"
    } else {
        Write-Host "⚠️  Regular: Posible optimización necesaria"
    }
}

Write-Host "`n💡 Nota: Esta prueba mide rendimiento del proxy SOCKS5 de IPv7"
Write-Host "Para pruebas más avanzadas, considera iperf con tunneling directo"