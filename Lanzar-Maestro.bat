@echo off
title IPv7-IEU (Nodo Maestro Gateway)
color 0A
echo ========================================================
echo        Iniciando IPv7-IEU (Gateway Maestro)
echo ========================================================
echo.
echo Este nodo actuara como salida a Internet para los satelites.
echo.
ipv7-c.exe -role master -port 7777 -api-port 7780 -tun=false
pause
