//go:build !windows && !linux && !darwin

package adapter

import "fmt"

// IEUInterface stub para plataformas sin soporte TUN nativo
type IEUInterface struct {
	IP   string
	Mask string
}

// NewIEUInterface retorna error informativo en plataformas no soportadas
func NewIEUInterface(name string, ip string, mask string) (*IEUInterface, error) {
	return nil, fmt.Errorf("adaptador virtual TUN no soportado en esta plataforma — usa -tun=false para operar en modo overlay puro")
}

func (i *IEUInterface) Read() ([]byte, error) {
	return nil, fmt.Errorf("TUN no disponible en esta plataforma")
}

func (i *IEUInterface) Write(data []byte) error {
	return fmt.Errorf("TUN no disponible en esta plataforma")
}

func (i *IEUInterface) Close() error {
	return nil
}
