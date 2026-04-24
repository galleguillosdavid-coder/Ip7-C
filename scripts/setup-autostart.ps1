# Script para configurar IPv7-IEU como servicio automático al inicio
# Requiere ejecutar como Administrador

Write-Host "🚀 Configurando IPv7-IEU para lanzamiento automático al encender la PC..."

# Ruta al ejecutable (ajusta si es diferente)
$exePath = "$PSScriptRoot\ipv7.exe"
$taskName = "IPv7-IEU-Master"

# Verificar que el ejecutable existe
if (!(Test-Path $exePath)) {
    Write-Host "❌ Error: No se encuentra ipv7.exe en $PSScriptRoot"
    exit 1
}

# Crear tarea programada para inicio del sistema
Write-Host "📅 Creando tarea programada '$taskName'..."

# Comando para ejecutar como Admin al inicio
$action = New-ScheduledTaskAction -Execute $exePath -Argument "--role master --port 7778 --api-port 7781 --tun=true"
$trigger = New-ScheduledTaskTrigger -AtStartup
$principal = New-ScheduledTaskPrincipal -UserId "SYSTEM" -LogonType ServiceAccount -RunLevel Highest
$settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -StartWhenAvailable

# Registrar la tarea
Register-ScheduledTask -TaskName $taskName -Action $action -Trigger $trigger -Principal $principal -Settings $settings -Description "IPv7-IEU Master Node - Inicia automáticamente al encender la PC"

Write-Host "✅ Tarea programada creada exitosamente."
Write-Host "ℹ️  La tarea '$taskName' se ejecutará como SYSTEM con privilegios elevados al iniciar Windows."
Write-Host "🔄 Para verificar: Ejecuta 'Get-ScheduledTask -TaskName $taskName' en PowerShell Admin."
Write-Host "🛑 Para detener: 'Stop-ScheduledTask -TaskName $taskName'"
Write-Host "🗑️  Para eliminar: 'Unregister-ScheduledTask -TaskName $taskName'"

Write-Host ""
Write-Host "⚠️  Nota: Si necesitas cambiar puertos o flags, edita la tarea en Task Scheduler o modifica este script."