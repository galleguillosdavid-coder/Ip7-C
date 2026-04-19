//go:build darwin

package adapter

import (
	"fmt"
	"net"
	"os/exec"

	"github.com/songgao/water"
)

// IEUInterface representa el adaptador virtual ieu0 (implementación macOS via utun)
type IEUInterface struct {
	iface *water.Interface
	IP    string
	Mask  string
}

// NewIEUInterface crea e inicializa una interfaz utun en macOS (requiere privilegios)
func NewIEUInterface(name string, ip string, mask string) (*IEUInterface, error) {
	cfg := water.Config{
		DeviceType: water.TUN,
	}
	// macOS asigna el nombre automáticamente (utun0, utun1, ...)
	iface, err := water.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("error al crear interfaz utun: %v", err)
	}

	// Calcular prefixlen desde máscara
	maskIP := net.ParseIP(mask).To4()
	prefixLen := 0
	if maskIP != nil {
		for _, b := range maskIP {
			for b != 0 {
				prefixLen += int(b & 1)
				b >>= 1
			}
		}
	} else {
		prefixLen = 24
	}

	// En macOS se usa ifconfig para configurar punto a punto
	// La dirección peer es la misma IP (loopback point-to-point)
	if err := exec.Command("ifconfig", iface.Name(), ip, ip, "up").Run(); err != nil {
		iface.Close()
		return nil, fmt.Errorf("error al configurar interfaz %s: %v", iface.Name(), err)
	}
	// Añadir ruta de red
	if err := exec.Command("route", "add", fmt.Sprintf("%s/%d", ip, prefixLen), ip).Run(); err != nil {
		// No fatal, puede ya existir la ruta
		fmt.Printf("⚠️ [macOS] No se pudo agregar ruta: %v\n", err)
	}

	return &IEUInterface{iface: iface, IP: ip, Mask: mask}, nil
}

// Read lee un paquete IP crudo de la interfaz utun
func (i *IEUInterface) Read() ([]byte, error) {
	// macOS TUN incluye 4-byte header de familia de protocolo
	buf := make([]byte, 65535)
	n, err := i.iface.Read(buf)
	if err != nil {
		return nil, err
	}
	// Omitir los 4 bytes del header de familia AF_INET
	if n > 4 {
		data := make([]byte, n-4)
		copy(data, buf[4:n])
		return data, nil
	}
	return buf[:n], nil
}

// Write inyecta un paquete IP crudo a la interfaz utun
func (i *IEUInterface) Write(data []byte) error {
	// Agregar 4-byte header AF_INET (0x00 0x00 0x00 0x02)
	packet := make([]byte, 4+len(data))
	packet[3] = 0x02 // AF_INET
	copy(packet[4:], data)
	_, err := i.iface.Write(packet)
	return err
}

// Close cierra la interfaz utun
func (i *IEUInterface) Close() error {
	return i.iface.Close()
}
