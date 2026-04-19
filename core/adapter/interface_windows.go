//go:build windows

package adapter

import (
	"fmt"
	"os/exec"

	"golang.org/x/sys/windows"
	"golang.zx2c4.com/wintun"
)

// IEUInterface representa el adaptador virtual ieu0 (implementación Windows via WinTun)
type IEUInterface struct {
	Adapter *wintun.Adapter
	Session wintun.Session
	IP      string
	Mask    string
}

// NewIEUInterface crea e inicializa el adaptador virtual WinTun (requiere Admin)
func NewIEUInterface(name string, ip string, mask string) (*IEUInterface, error) {
	adapter, err := wintun.CreateAdapter(name, "Wintun", nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear adaptador (asegúrate de correr como Admin): %v", err)
	}

	cmd := exec.Command("netsh", "interface", "ip", "set", "address", name, "static", ip, mask)
	if err := cmd.Run(); err != nil {
		adapter.Close()
		return nil, fmt.Errorf("error al configurar IP en %s: %v", name, err)
	}

	session, err := adapter.StartSession(0x400000) // 4MB ring buffer
	if err != nil {
		adapter.Close()
		return nil, fmt.Errorf("error al iniciar sesión WinTun: %v", err)
	}

	return &IEUInterface{
		Adapter: adapter,
		Session: session,
		IP:      ip,
		Mask:    mask,
	}, nil
}

// Read lee un paquete IP crudo bloqueando hasta que haya datos disponibles
func (i *IEUInterface) Read() ([]byte, error) {
	for {
		packet, err := i.Session.ReceivePacket()
		if err == nil {
			defer i.Session.ReleaseReceivePacket(packet)
			data := make([]byte, len(packet))
			copy(data, packet)
			return data, nil
		}
		if err == windows.ERROR_NO_MORE_ITEMS {
			windows.WaitForSingleObject(i.Session.ReadWaitEvent(), windows.INFINITE)
			continue
		}
		return nil, err
	}
}

// Write escribe un paquete IP crudo hacia el sistema operativo
func (i *IEUInterface) Write(data []byte) error {
	packet, err := i.Session.AllocateSendPacket(len(data))
	if err != nil {
		return err
	}
	copy(packet, data)
	i.Session.SendPacket(packet)
	return nil
}

// Close libera los recursos del adaptador
func (i *IEUInterface) Close() error {
	i.Session.End()
	if i.Adapter != nil {
		return i.Adapter.Close()
	}
	return nil
}
