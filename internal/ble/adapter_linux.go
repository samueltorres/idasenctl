package ble

import (
	"tinygo.org/x/bluetooth"
)

func bluetoothAddress(address string) (*bluetooth.Address, error) {
	addr, err := bluetooth.ParseMAC(address)
	if err != nil {
		return nil, err
	}

	return &bluetooth.Address{MACAddress: bluetooth.MACAddress{MAC: addr}}, nil
}
