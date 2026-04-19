Para lanzar esta aplicación y que todo funcione como lo hemos configurado (especialmente el acceso a Internet y el Proxy), debes abrir dos terminales (PowerShell) y ejecutar los siguientes comandos:

1. El Nodo Maestro (Es el que tiene la salida a Internet)
Este nodo recibe las peticiones de los satélites y las saca al mundo real.

powershell
.\ipv7-node.exe --role master --port 7777 --api-port 7780 --tun=false
Enrutador Global: Se activa automáticamente en el puerto TCP 9777.
2. El Nodo Satélite (Es el que usas para navegar)
Este nodo crea el punto de acceso para tu navegador o móvil. Conéctalo al Maestro usando su IP (usa 127.0.0.1 si estás en la misma PC).

powershell
.\ipv7.exe --role node --port 7778 --api-port 7781 --tun=false --bootstrap 127.0.0.1:8777 --remote 127.0.0.1
Proxy SOCKS5/HTTP: Se activa en el puerto 1080 de esta máquina.
3. Cómo conectar tus dispositivos
Una vez que ambos nodos digan "¡Éxito!" y veas los corazones latiendo (Keep-alive), configura tu conexión así:

En esta misma PC: Configura el proxy de Windows o tu navegador en 127.0.0.1 puerto 1080 (tipo SOCKS5 o HTTP, ambos funcionan gracias al parche que puse).
En tu Móvil u otro PC: Usa la IP local de esta computadora (ej. 192.168.1.198) y el puerto 1080.
Nota importante: Como usamos --tun=false, la aplicación no necesita permisos de Administrador y no crea una tarjeta de red virtual "real", sino que funciona como un túnel de aplicación (Proxy) extremadamente rápido. Todo lo que pase por ese puerto 1080 volará a través de IPv7.

¿Quieres que proceda ahora con el commit y envío a GitHub de todo este nuevo sistema de Proxy que programamos?

si