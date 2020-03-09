package shelly

import (
	"github.com/reef-pi/hal"
	"testing"
)

func TestShelly25(t *testing.T) {
	d, err := NewShelly25("127.0.0.1", true)
	if err != nil {
		t.Error(err)
	}
	if d.Metadata().Name == "" {
		t.Error("HAL metadata should not have empty name")
	}

	d1 := d.(hal.DigitalOutputDriver)

	if len(d1.DigitalOutputPins()) != 2 {
		t.Error("Expected exactly two output pin")
	}
	pin, err := d1.DigitalOutputPin(0)
	if err != nil {
		t.Error(err)
	}
	if pin.LastState() != false {
		t.Error("Expected initial state to be false")
	}
}
