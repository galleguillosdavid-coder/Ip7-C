# Script para lanzar múltiples instancias IPv7-IEU con sub-puertos
# Uso: .\lanzar-instancia.ps1 -SubPort 2 -Role node

param(
    [Parameter(Mandatory=$true)]
    [int]$SubPort,
    
    [Parameter(Mandatory=$false)]
    [string]$Role = "node",
    
    [Parameter(Mandatory=$false)]
    [int]$Port = 7780 + $SubPort,
    
    [Parameter(Mandatory=$false)]
    [int]$ApiPort = 7783 + $SubPort
)

Write-Host "🚀 Lanzando IPv7-IEU Instancia $SubPort (Rol: $Role)"
Write-Host "=============================================="

$baseArgs = @(
    "--role", $Role,
    "--port", $Port,
    "--api-port", $ApiPort,
    "--tun=true",
    "--sub-port", $SubPort
)

if ($Role -eq "node") {
    $baseArgs += @("--bootstrap", "127.0.0.1:8778", "--remote", "127.0.0.1", "--remote-port", "7778")
}

Write-Host "Comando: .\ipv7.exe $($baseArgs -join ' ')"
Start-Process -FilePath ".\ipv7.exe" -ArgumentList $baseArgs -Verb RunAs -WindowStyle Hidden

Write-Host "✅ Instancia $SubPort lanzada en background"