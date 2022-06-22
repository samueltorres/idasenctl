package idasen

import (
	"context"
	"encoding/binary"
	"errors"
	"math"

	"github.com/samueltorres/idasenctl/internal/ble"
)

var (
	IDASEN_MIN_HEIGHT                   = 0.62
	IDASEN_MAX_HEIGHT                   = 1.27
	IDASEN_UUID_HEIGHT                  = "99fa0021-338a-1024-8a49-009c0215f78a"
	IDASEN_UUID_COMMAND                 = "99fa0002-338a-1024-8a49-009c0215f78a"
	IDASEN_UUID_REFERENCE_INPUT         = "99fa0031-338a-1024-8a49-009c0215f78a"
	IDASEN_COMMAND_REFERENCE_INPUT_STOP = []byte{0x01, 0x80}
	IDASEN_COMMAND_UP                   = []byte{0x47, 0x00}
	IDASEN_COMMAND_DOWN                 = []byte{0x46, 0x00}
	IDASEN_COMMAND_STOP                 = []byte{0xFF, 0x00}

	ErrHeightBiggerThanMax  = errors.New("height is bigger than the max height")
	ErrHeightSmallerThanMin = errors.New("height is smaller than the min height")
)

type Controller struct {
	adaptor *ble.Adapter
}

func NewController(deskAddress string) (*Controller, error) {
	bleAdaptor, err := ble.NewAdapter(deskAddress)
	if err != nil {
		return nil, err
	}

	return &Controller{
		adaptor: bleAdaptor,
	}, nil
}

func (c *Controller) MoveTo(ctx context.Context, desiredHeight float32, updates chan<- float32) error {
	if desiredHeight > float32(IDASEN_MAX_HEIGHT) {
		return ErrHeightBiggerThanMax
	}

	if desiredHeight < float32(IDASEN_MIN_HEIGHT) {
		return ErrHeightBiggerThanMax
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		currentHeight, err := c.GetCurrentHeight()
		if err != nil {
			return err
		}

		if updates != nil {
			updates <- currentHeight
		}

		if math.Abs(float64(desiredHeight-currentHeight)) < 0.005 {
			err := c.stop()
			if err != nil {
				return err
			}
			break
		}

		if desiredHeight <= currentHeight {
			err := c.moveDown()
			if err != nil {
				return err
			}
		} else {
			err := c.moveUp()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Controller) GetCurrentHeight() (float32, error) {
	b, err := c.adaptor.ReadCharacteristic(IDASEN_UUID_HEIGHT)
	if err != nil {
		return 0, err
	}
	raw := binary.LittleEndian.Uint16(b[0:2])
	return float32(float32(raw)/10000) + float32(IDASEN_MIN_HEIGHT), err
}

func (c *Controller) moveUp() error {
	err := c.adaptor.WriteCharacteristic(IDASEN_UUID_COMMAND, IDASEN_COMMAND_UP)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) moveDown() error {
	err := c.adaptor.WriteCharacteristic(IDASEN_UUID_COMMAND, IDASEN_COMMAND_DOWN)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) stop() error {
	err := c.adaptor.WriteCharacteristic(IDASEN_UUID_REFERENCE_INPUT, IDASEN_COMMAND_STOP)
	if err != nil {
		return err
	}
	return nil
}
