package idasen

import (
	"context"
	"strings"

	"tinygo.org/x/bluetooth"
)

type Scanner struct {
	adapter *bluetooth.Adapter
}

type DeviceInfo struct {
	Name    string
	Address string
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

func (s *Scanner) Scan(ctx context.Context, dic chan<- DeviceInfo) {
	desksSeen := make(map[string]struct{})
	s.adapter.Scan(func(a *bluetooth.Adapter, device bluetooth.ScanResult) {
		if strings.HasPrefix(device.LocalName(), "Desk") {
			if _, ok := desksSeen[device.LocalName()]; !ok {
				dic <- DeviceInfo{Name: device.LocalName(), Address: device.Address.String()}
			}
			desksSeen[device.LocalName()] = struct{}{}

			select {
			case <-ctx.Done():
				a.StopScan()
			default:
			}
		}
	})
}
