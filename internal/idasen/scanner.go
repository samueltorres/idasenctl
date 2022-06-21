package idasen

import (
	"fmt"
	"strings"

	"tinygo.org/x/bluetooth"
)

type Scanner struct {
	adapter *bluetooth.Adapter
}

func NewScanner() (*Scanner, error) {
	adapter := bluetooth.DefaultAdapter
	err := adapter.Enable()
	if err != nil {
		return nil, err
	}
	return &Scanner{
		adapter: adapter,
	}, nil
}

func (s *Scanner) Scan() {
	fmt.Println("Scanning")
	s.adapter.Scan(func(a *bluetooth.Adapter, device bluetooth.ScanResult) {
		if strings.HasPrefix(device.LocalName(), "Desk") {
			fmt.Println("Found Desk", device.LocalName(), device.Address.String())
			a.StopScan()
		}
	})
}
