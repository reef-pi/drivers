package file

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"encoding/json"
	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
)

type Config struct {
	Address string `json:"address"`
}

type analog struct {
	path       string
	meta       hal.Metadata
	calibrator hal.Calibrator
}

func HalAnalogAdapter(c []byte, _ i2c.Bus) (hal.Driver, error) {
	var config Config
	if err := json.Unmarshal(c, &config); err != nil {
		return nil, err
	}
	return NewAnalog(config.Address)
}

func NewAnalog(p string) (*analog, error) {
	c, err := hal.CalibratorFactory([]hal.Measurement{})
	if err != nil {
		return nil, err
	}
	return &analog{
		path:       p,
		calibrator: c,
		meta: hal.Metadata{
			Name:         "analog-file",
			Description:  "A simple file based analog hal driver",
			Capabilities: []hal.Capability{hal.AnalogInput},
		},
	}, nil
}

func (f *analog) Metadata() hal.Metadata {
	return f.meta
}

func (f *analog) Close() error {
	return nil
}

func (f *analog) Name() string {
	return f.path
}

func (f *analog) Number() int {
	return 0
}
func (f *analog) Read() (float64, error) {
	data, err := ioutil.ReadFile(f.path)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
}

func (f *analog) Measure() (float64, error) {
	v, err := f.Read()
	if err != nil {
		return 0, err
	}
	if f.calibrator == nil {
		return 0, fmt.Errorf("Not calibrated")
	}
	return f.calibrator.Calibrate(v), nil
}
func (f *analog) Calibrate(points []hal.Measurement) error {
	cal, err := hal.CalibratorFactory(points)
	if err != nil {
		return err
	}
	f.calibrator = cal
	return nil
}

func (f *analog) AnalogInputPins() []hal.AnalogInputPin {
	return []hal.AnalogInputPin{f}
}

func (f *analog) AnalogInputPin(_ int) (hal.AnalogInputPin, error) {
	return f, nil
}

func (f *analog) Pins(cap hal.Capability) ([]hal.Pin, error) {
	if cap == hal.AnalogInput {
		return []hal.Pin{f}, nil
	}
	return nil, fmt.Errorf("unsupported capability:%s", cap.String())
}
