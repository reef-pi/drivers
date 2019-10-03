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

type digital struct {
	path      string
	meta      hal.Metadata
	lastState bool
}

func HalDigitalAdapter(c []byte, _ i2c.Bus) (hal.Driver, error) {
	var config Config
	if err := json.Unmarshal(c, &config); err != nil {
		return nil, err
	}
	return NewDigital(config.Address), nil
}

func NewDigital(p string) *digital {
	return &digital{
		path: p,
		meta: hal.Metadata{
			Name:         "digital-file",
			Description:  "A simple file based digital hal driver",
			Capabilities: []hal.Capability{hal.DigitalInput, hal.DigitalOutput, hal.PWM},
		},
	}
}

func (f *digital) Metadata() hal.Metadata {
	return f.meta
}

func (d *digital) Number() int {
	return 0
}
func (f *digital) Close() error {
	return nil
}

func (f *digital) Name() string {
	return f.path
}

func (f *digital) Read() (bool, error) {
	data, err := ioutil.ReadFile(f.path)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(data)) == "1", nil

}

func (f *digital) LastState() bool {
	return f.lastState
}

func (f *digital) Write(b bool) error {
	f.lastState = b
	if b {
		return ioutil.WriteFile(f.path, []byte("1"), 0644)
	}
	return ioutil.WriteFile(f.path, []byte("0"), 0644)

}
func (f *digital) Set(v float64) error {
	return ioutil.WriteFile(f.path, []byte(strconv.FormatFloat(v, 'f', -1, 64)), 0644)
}

func (f *digital) DigitalInputPins() []hal.DigitalInputPin {
	return []hal.DigitalInputPin{f}
}

func (f *digital) DigitalInputPin(_ int) (hal.DigitalInputPin, error) {
	return f, nil
}

func (f *digital) DigitalOutputPins() []hal.DigitalOutputPin {
	return []hal.DigitalOutputPin{f}
}

func (f *digital) DigitalOutputPin(_ int) (hal.DigitalOutputPin, error) {
	return f, nil
}
func (f *digital) PWMChannels() []hal.PWMChannel {
	return []hal.PWMChannel{f}
}

func (f *digital) PWMChannel(_ int) (hal.PWMChannel, error) {
	return f, nil
}
func (f *digital) Pins(cap hal.Capability) ([]hal.Pin, error) {
	switch cap {
	case hal.DigitalInput, hal.DigitalOutput, hal.PWM:
		return []hal.Pin{f}, nil
	default:
		return nil, fmt.Errorf("unsupported capability:%s", cap.String())
	}
}
