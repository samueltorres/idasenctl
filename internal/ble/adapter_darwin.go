package ble

import (
	"tinygo.org/x/bluetooth"
)

func bluetoothAddress(address string) (*bluetooth.Address, error) {
	deviceUUID, err := bluetooth.ParseUUID(address)
	if err != nil {
		return nil, err
	}

	return &bluetooth.Address{UUID: deviceUUID}, nil
}
