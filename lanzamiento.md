Para lanzar esta aplicación de manera rápida y sin complicaciones de terminal, he creado para ti **2 accesos directos (.bat)** en la carpeta del proyecto. Solamente necesitas hacerles doble clic:

### 1. El Nodo Maestro `Lanzar-Maestro.bat`
Este nodo recibe las peticiones de los satélites y las saca al mundo real.
*(Enrutador Global activado automáticamente en el puerto TCP 9777).*

### 2. El Nodo Satélite `Lanzar-Satelite.bat`
Este nodo crea el punto de acceso Proxy SOCKS5 (puerto 1080) para tu navegador o móvil. Conéctalo y úsalo para navegar.
3. Cómo conectar tus dispositivos
Una vez que ambos nodos digan "¡Éxito!" y veas los corazones latiendo (Keep-alive), configura tu conexión así:

En esta misma PC: Configura el proxy de Windows o tu navegador en 127.0.0.1 puerto 1080 (tipo SOCKS5 o HTTP, ambos funcionan gracias al parche que puse).
En tu Móvil u otro PC: Usa la IP local de esta computadora (ej. 192.168.1.198) y el puerto 1080.
Nota importante: Como usamos --tun=false, la aplicación no necesita permisos de Administrador y no crea una tarjeta de red virtual "real", sino que funciona como un túnel de aplicación (Proxy) extremadamente rápido. Todo lo que pase por ese puerto 1080 volará a través de IPv7.

¿Quieres que proceda ahora con el commit y envío a GitHub de todo este nuevo sistema de Proxy que programamos?

si