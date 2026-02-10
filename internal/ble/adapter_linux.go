package ble

import (
	"tinygo.org/x/bluetooth"
)

func bluetoothAddress(address string) (bluetooth.Address, error) {
	addr, err := bluetooth.ParseMAC(address)
	if err != nil {
		return bluetooth.Address{}, err
	}

	return bluetooth.Address{MACAddress: bluetooth.MACAddress{MAC: addr}}, nil
}
