package ble

import (
	"errors"
	"time"

	"tinygo.org/x/bluetooth"
)

var (
	ErrCharacteristicNotExists = errors.New("characteristic does not exist")
)

type Adapter struct {
	device *bluetooth.Device
}

func NewAdapter(address string) (*Adapter, error) {
	adapter := bluetooth.DefaultAdapter
	adapter.Enable()
	deviceUUID, err := bluetooth.ParseUUID(address)
	if err != nil {
		return nil, err
	}

	device, err := adapter.Connect(
		bluetooth.Address{UUID: deviceUUID},
		bluetooth.ConnectionParams{
			ConnectionTimeout: bluetooth.NewDuration(30 * time.Second),
		})

	if err != nil {
		return nil, err
	}

	bleAdapter := &Adapter{
		device: device,
	}

	return bleAdapter, nil
}

func (a *Adapter) ReadCharacteristic(cUUID string) ([]byte, error) {
	c, err := a.getCharacteristic(cUUID)
	if err != nil {
		return nil, ErrCharacteristicNotExists
	}

	b := make([]byte, 16)
	_, err = c.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (a *Adapter) WriteCharacteristic(cUUID string, data []byte) error {
	c, err := a.getCharacteristic(cUUID)
	if err != nil {
		return ErrCharacteristicNotExists
	}
	if c == nil {
		return ErrCharacteristicNotExists
	}
	_, err = c.WriteWithoutResponse(data)
	if err != nil {
		return err
	}

	return nil
}

func (a *Adapter) getCharacteristic(cUUID string) (*bluetooth.DeviceCharacteristic, error) {
	ds, err := a.device.DiscoverServices(nil)
	if err != nil {
		return nil, err
	}

	for _, s := range ds {
		c, err := s.DiscoverCharacteristics(nil)
		if err != nil {
			return nil, err
		}

		for _, cc := range c {
			if cc.UUID().String() == cUUID {
				return &cc, nil
			}
		}
	}
	return nil, ErrCharacteristicNotExists
}
