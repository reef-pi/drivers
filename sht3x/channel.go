package sht3x

import (
	"fmt"
	"github.com/reef-pi/hal"
)

type channel struct {
	calibrator hal.Calibrator
	d          *SHT31D
	number     int
}

func newChannel(d *SHT31D, i int) (hal.AnalogInputPin, error) {
	c, err := hal.CalibratorFactory([]hal.Measurement{})
	if err != nil {
		return nil, err
	}
	return &channel{
		calibrator: c,
		number:     i,
		d:          d,
	}, nil
}

func (c *channel) Name() string {
	switch c.number {
	case 0:
		return "temperature"
	case 1:
		return "humidity"
	default:
		return "unknown"
	}
}
func (c *channel) Number() int {
	return c.number
}

func (c *channel) Calibrate(points []hal.Measurement) error {
	cal, err := hal.CalibratorFactory(points)
	if err != nil {
		return err
	}
	c.calibrator = cal
	return nil
}

func (c *channel) Value() (float64, error) {
	switch c.number {
	case 0:
		return c.d.Temperature()
	case 1:
		return c.d.Humidity()
	default:
		return 0, nil
	}
}

func (c *channel) Measure() (float64, error) {
	v, err := c.Value()
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
