package ads1x15

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
)

type channel struct {
	bus        i2c.Bus
	address    byte
	pinAddress uint16
	channel    int
	gainConfig uint16
	calibrator hal.Calibrator
	shift      int
	delay      time.Duration
}

func newChannel(b i2c.Bus, address byte, channelNum int, pinAddress uint16, gain uint16, shift int, delay time.Duration) (*channel, error) {
	c, err := hal.CalibratorFactory([]hal.Measurement{})
	if err != nil {
		return nil, err
	}

	return &channel{
		bus:        b,
		address:    address,
		pinAddress: pinAddress,
		channel:    channelNum,
		gainConfig: gain,
		calibrator: c,
		shift:      shift,
		delay:      delay,
	}, nil

}

func (c *channel) Name() string {
	return fmt.Sprintf("%d", c.channel)
}

func (c *channel) Number() int {
	return c.channel
}

func (c *channel) Calibrate(points []hal.Measurement) error {
	cal, err := hal.CalibratorFactory(points)
	if err != nil {
		return err
	}
	c.calibrator = cal
	return nil
}

func (c *channel) Read() (float64, error) {

	var config uint16 = configOsSingle |
		configModeSingle |
		configDataRate1600 |
		configComparatorModeTraditional |
		configComparitorNonLatching |
		configComparitorPolarityActiveLow |
		configComparitorQueueNone |
		c.pinAddress |
		c.gainConfig

	configBytes := make([]byte, 2)

	binary.BigEndian.PutUint16(configBytes, uint16(config))

	var v float64
	var e error

	for attempt := 0; attempt <= 3; attempt++ {
		v, e = c.performConversion(configBytes[:])

		if e == nil {
			return v, e
		}
	}

	return v, e
}

func (c *channel) Measure() (float64, error) {
	v, err := c.Read()
	if err != nil {
		return 0, err
	}
	if c.calibrator == nil {
		return 0, fmt.Errorf("Not calibrated")
	}
	return c.calibrator.Calibrate(v), nil
}

func (c *channel) Close() error {
	return nil
}

func (c *channel) performConversion(configBytes []byte) (float64, error) {

	b := make([]byte, 2)

	if err := c.bus.WriteToReg(c.address, 0x01, configBytes[:]); err != nil {
		return 0, err
	}

	time.Sleep(c.delay)

	if err := c.bus.ReadFromReg(c.address, 0x01, b[:]); err != nil {
		return 0, err
	}

	if configBytes[0] != b[0] || configBytes[1] != b[1] {
		return 0, errors.New("config mismatch")
	}

	if err := c.bus.ReadFromReg(c.address, 0x00, b[:]); err != nil {
		return 0, err
	}

	v := int16(int16(b[0])<<8|int16(b[1])) >> c.shift

	return float64(v), nil
}
