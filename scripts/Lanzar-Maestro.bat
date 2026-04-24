@echo off
title IPv7-IEU (Nodo Maestro Gateway)
color 0A
echo ========================================================
echo        Iniciando IPv7-IEU (Gateway Maestro)
echo ========================================================
echo.
echo Este nodo actuara como salida a Internet para los satelites.
echo Modo Proxy (sin TUN, no requiere Admin).
echo.
ipv7.exe --role master --port 7778 --api-port 7781 --tun=false --sub-port 0
pause
