package tasmota

import (
	"github.com/reef-pi/hal"
	"testing"
)

func TestHttpDriver(t *testing.T) {

	f := HttpDriverFactory()

	params := map[string]interface{}{
		"Domain Or Address": "192.168.1.46",
	}

	d, err := f.NewDriver(params, nil)
	if err != nil {
		t.Fatal(err)
	}

	meta := d.Metadata()
	if len(meta.Capabilities) != 2 {
		t.Error("Expected 1 capabilities, found:", len(meta.Capabilities))
	}

	o, ok := d.(hal.DigitalOutputDriver)
	if !ok {
		t.Error("Failed to type driver to Digital output driver")
	}

	if len(o.DigitalOutputPins()) != 1 {
		t.Error("Expected a single digital output pwm pin, found:", len(o.DigitalOutputPins()))
	}

	_, err = o.DigitalOutputPin(0)
	if err != nil {
		t.Error("Expected a digital output pin")
	}

	pwm, ok := d.(hal.PWMDriver)
	if !ok {
		t.Error("Failed to type driver to PWM driver")
	}

	if len(pwm.PWMChannels()) != 1 {
		t.Error("Expected a single pwm channel, found:", len(pwm.PWMChannels()))
	}

	_, err = pwm.PWMChannel(0)
	if err != nil {
		t.Error("Expected a pwm pin")
	}

}
