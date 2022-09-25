package esp32

import (
	"fmt"
	"github.com/reef-pi/hal"
	"net/http"
	"time"
)

const (
	_driverName = "esp32"
	_timeout    = 3 * time.Second
)

type HTTPClient func(*http.Request) (*http.Response, error)

type driver struct {
	meta    hal.Metadata
	address string
	pins    map[hal.Capability][]int
	client  HTTPClient
}

func (d *driver) Close() error {
	return nil
}

func (m *driver) Metadata() hal.Metadata {
	return m.meta
}

func (m *driver) Name() string {
	return _driverName
}

func (d *driver) Pins(cap hal.Capability) ([]hal.Pin, error) {
	ps, ok := d.pins[cap]
	if !ok {
		return nil, fmt.Errorf("unsupported capability:%s", cap.String())
	}
	var pins []hal.Pin
	for _, p := range ps {
		pins = append(pins, d.halPin(cap, p))
	}
	return pins, nil
}

func (d *driver) PWMChannels() []hal.PWMChannel {
	var channels []hal.PWMChannel
	for _, p := range d.pins[hal.PWM] {
		channels = append(channels, d.halPin(hal.PWM, p))
	}
	return channels
}

func (d *driver) halPin(c hal.Capability, p int) *pin {
	return &pin{
		address: d.address,
		number:  p,
		cap:     c,
		client:  d.client,
	}
}

func (d *driver) PWMChannel(i int) (hal.PWMChannel, error) {
	for _, p := range d.pins[hal.PWM] {
		if p == i {
			return d.halPin(hal.PWM, p), nil
		}
	}
	return nil, fmt.Errorf("no pwm channels for pin %d found", i)
}

func (d *driver) DigitalOutputPins() []hal.DigitalOutputPin {
	var pins []hal.DigitalOutputPin
	for _, p := range d.pins[hal.DigitalOutput] {
		pins = append(pins, d.halPin(hal.DigitalOutput, p))
	}
	return pins
}

func (d *driver) DigitalOutputPin(i int) (hal.DigitalOutputPin, error) {
	for _, p := range d.pins[hal.DigitalOutput] {
		if p == i {
			return d.halPin(hal.DigitalOutput, p), nil
		}
	}
	return nil, fmt.Errorf("no pwm channels for pin %d found", i)
}

func (d *driver) DigitalInputPins() []hal.DigitalInputPin {

	var pins []hal.DigitalInputPin
	for _, p := range d.pins[hal.DigitalInput] {
		pins = append(pins, d.halPin(hal.DigitalInput, p))
	}
	return pins
}

func (d *driver) DigitalInputPin(i int) (hal.DigitalInputPin, error) {
	for _, p := range d.pins[hal.DigitalInput] {
		if p == i {
			return d.halPin(hal.DigitalInput, p), nil
		}
	}
	return nil, fmt.Errorf("no pwm channels for pin %d found", i)
}
func (d *driver) AnalogInputPins() []hal.AnalogInputPin {

	var pins []hal.AnalogInputPin
	for _, p := range d.pins[hal.AnalogInput] {
		pins = append(pins, d.halPin(hal.AnalogInput, p))
	}
	return pins
}

func (d *driver) AnalogInputPin(i int) (hal.AnalogInputPin, error) {
	for _, p := range d.pins[hal.AnalogInput] {
		if p == i {
			return d.halPin(hal.AnalogInput, p), nil
		}
	}
	return nil, fmt.Errorf("no analog input pin %d found", i)
}
