//go:build linux

package adapter

import (
	"fmt"
	"net"
	"os/exec"

	"github.com/songgao/water"
)

// IEUInterface representa el adaptador virtual ieu0 (implementación Linux via TUN)
type IEUInterface struct {
	iface *water.Interface
	IP    string
	Mask  string
}

// NewIEUInterface crea e inicializa un adaptador TUN en Linux (requiere root o CAP_NET_ADMIN)
func NewIEUInterface(name string, ip string, mask string) (*IEUInterface, error) {
	cfg := water.Config{
		DeviceType: water.TUN,
	}
	cfg.Name = name

	iface, err := water.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("error al crear interfaz TUN (¿necesitas root?): %v", err)
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
		prefixLen = 24 // default /24
	}

	// Asignar IP y levantar interfaz
	if err := exec.Command("ip", "addr", "add", fmt.Sprintf("%s/%d", ip, prefixLen), "dev", iface.Name()).Run(); err != nil {
		iface.Close()
		return nil, fmt.Errorf("error al asignar IP %s: %v", ip, err)
	}
	if err := exec.Command("ip", "link", "set", "dev", iface.Name(), "up").Run(); err != nil {
		iface.Close()
		return nil, fmt.Errorf("error al levantar interfaz %s: %v", iface.Name(), err)
	}

	return &IEUInterface{iface: iface, IP: ip, Mask: mask}, nil
}

// Read lee un paquete IP crudo de la interfaz TUN
func (i *IEUInterface) Read() ([]byte, error) {
	buf := make([]byte, 65535)
	n, err := i.iface.Read(buf)
	if err != nil {
		return nil, err
	}
	data := make([]byte, n)
	copy(data, buf[:n])
	return data, nil
}

// Write inyecta un paquete IP crudo al kernel de Linux
func (i *IEUInterface) Write(data []byte) error {
	_, err := i.iface.Write(data)
	return err
}

// Close cierra y elimina la interfaz TUN
func (i *IEUInterface) Close() error {
	return i.iface.Close()
}
