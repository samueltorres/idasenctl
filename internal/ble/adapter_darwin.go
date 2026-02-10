package ble

import (
	"tinygo.org/x/bluetooth"
)

func bluetoothAddress(address string) (bluetooth.Address, error) {
	deviceUUID, err := bluetooth.ParseUUID(address)
	if err != nil {
		return bluetooth.Address{}, err
	}

	return bluetooth.Address{UUID: deviceUUID}, nil
}
