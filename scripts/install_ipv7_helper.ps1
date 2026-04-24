#Requires -RunAsAdministrator

$src = Join-Path $PSScriptRoot '..\ipv7.exe' | Resolve-Path -ErrorAction Stop
$dstDir = Join-Path $env:LOCALAPPDATA 'Programs\IPv7'
$dst = Join-Path $dstDir 'ipv7.exe'
$bat = Join-Path $env:APPDATA 'Microsoft\Windows\Start Menu\Programs\Startup\Start-IPv7-IEU.bat'

if (-not (Test-Path $src)) {
    Write-Error "Source file missing: $src"
    exit 1
}

New-Item -ItemType Directory -Force -Path $dstDir | Out-Null
Copy-Item -Path $src -Destination $dst -Force

$lines = @(
    '@echo off',
    "cd /d `"$dstDir`"",
    "start `"IPv7-IEU`" `"$dst`" --role master --port 7778 --api-port 7781 --tun=true"
)
Set-Content -Path $bat -Value $lines -Encoding ASCII

Write-Output "Installed: $dst"
Write-Output "Startup: $bat"
Get-Item $dst | Select-Object FullName,Length,LastWriteTime
Get-Content $bat
