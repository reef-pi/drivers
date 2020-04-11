package sht3x

import (
	"testing"

	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
)

var params = map[string]interface{}{
	"Address": 0x44,
}

func TestDriver(t *testing.T) {
	bus := i2c.MockBus()
	//bus.Bytes = make([]byte, 2)

	f := Factory()
	_, err := f.NewDriver(nil, bus)

	if err == nil {
		t.Error("Adapter creation should fail when json config is invalid")
	}

	driver, err := f.NewDriver(params, bus)

	if err != nil {
		t.Error(err)
	}
	if driver.Metadata().Name != "sht31d" {
		t.Error("Unexpected name")
	}
	if !driver.Metadata().HasCapability(hal.AnalogInput) {
		t.Error("Analog input Capability should exist")
	}
	if driver.Metadata().HasCapability(hal.DigitalInput) {
		t.Error("Digital Input Capability should not exist")
	}

	d := driver.(hal.AnalogInputDriver)

	if len(d.AnalogInputPins()) != 2 {
		t.Error("Expected only one channel")
	}
	if _, err := d.AnalogInputPin(3); err == nil {
		t.Error("Expected error for invalid channel name")
	}

	ch, err := d.AnalogInputPin(0)
	if err != nil {
		t.Error(err)
	}
	if ch.Name() != "temperature" {
		t.Error("Unexpected channel name")
	}
	bus.Bytes = []byte{0x60, 0xc4, 0x57, 0x7f, 0x15, 0x95}
	v, err := ch.Read()
	if err != nil {
		t.Error(err)
	}
	if v != 21.149385824368665 {
		t.Error("Unexepected value:", v)
	}
	if err := d.Close(); err != nil {
		t.Error(err)
	}
}
