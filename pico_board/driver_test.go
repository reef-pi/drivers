package pico_board

import (
	"testing"

	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
)

var params = map[string]interface{}{
	"Address": 72,
}

func TestPhBoardDriver(t *testing.T) {
	bus := i2c.MockBus()
	bus.Bytes = make([]byte, 2)

	f := Factory()
	_, err := f.NewDriver(nil, bus)

	if err == nil {
		t.Error("Adapter creation should fail when json config is invalid")
	}

	driver, err := f.NewDriver(params, bus)

	if err != nil {
		t.Error(err)
	}
	if driver.Metadata().Name != "pico-board" {
		t.Error("Unexpected name")
	}
	if !driver.Metadata().HasCapability(hal.AnalogInput) {
		t.Error("analog input cpability should exist")
	}
	if driver.Metadata().HasCapability(hal.DigitalInput) {
		t.Error("Digital input Capability should not exist")
	}

	d := driver.(hal.AnalogInputDriver)

	if len(d.AnalogInputPins()) != 1 {
		t.Error("Expected only one channel")
	}
	if _, err := d.AnalogInputPin(1); err == nil {
		t.Error("Expected error for invalid channel name")
	}

	ch, err := d.AnalogInputPin(0)
	if err != nil {
		t.Error(err)
	}
	if ch.Name() != "0" {
		t.Error("Unexpected channel name")
	}
	v, err := ch.Value()
	if err != nil {
		t.Error(err)
	}
	if v != 0 {
		t.Error("Unexepected value")
	}
	if err := d.Close(); err != nil {
		t.Error(err)
	}
}

func TestFloatAddress(t *testing.T) {
	bus := i2c.MockBus()
	bus.Bytes = make([]byte, 2)

	var floatAddress float64
	floatAddress = 72
	var floatparams = map[string]interface{}{
		"Address": floatAddress,
	}

	f := Factory()
	_, err := f.NewDriver(floatparams, bus)

	if err != nil {
		t.Error("pico_board should convert address to int before casting to byte")
	}
}
