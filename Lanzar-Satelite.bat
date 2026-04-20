@echo off
title IPv7-IEU (Nodo Satelite Cliente)
color 0B
echo ========================================================
echo        Iniciando IPv7-IEU (Satelite Proxy)
echo ========================================================
echo.
echo Configurando Proxy de Windows 11 (127.0.0.1:1080)...
reg add "HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings" /v ProxyServer /t REG_SZ /d "127.0.0.1:1080" /f >nul
reg add "HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings" /v ProxyEnable /t REG_DWORD /d 1 /f >nul

echo Conectando al Maestro local y abriendo Proxy SOCKS5 en puerto 1080...
echo.
ipv7-c.exe -role node -port 7778 -api-port 7781 -tun=false -bootstrap 127.0.0.1:8777 -remote 127.0.0.1

echo.
echo Desactivando Proxy de Windows (Restaurando conexion normal)...
reg add "HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings" /v ProxyEnable /t REG_DWORD /d 0 /f >nul
pause
