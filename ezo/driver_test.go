package ezo

import (
	"testing"

	"github.com/reef-pi/hal"

	"github.com/reef-pi/rpi/i2c"
)

func TestEZO(t *testing.T) {
	factory := Factory()
	params := map[string]interface{}{
		"Address": 0x93,
	}

	bus := i2c.MockBus()

	driver, err := factory.NewDriver(params, bus)

	if err != nil {
		t.Error("Unable to crate EZO driver")
	}

	e, ok := driver.(*AtlasEZO)
	if !ok {
		t.Error("Unable to convert driver to AtlasEZO")
	}

	e.delay = 0
	bus.Bytes = append([]byte{1}, []byte("9.65")...)
	if _, err := e.Read(); err != nil {
		t.Error(err)
	}
	bus.Bytes = append([]byte{0}, []byte("9.65")...)
	if _, err := e.Read(); err == nil {
		t.Error("Values starting with 0 should fail")
	}

	bus.Bytes = append([]byte{1}, []byte("L,1")...)
	on, err := e.LedState()
	if err != nil {
		t.Error(err)
	}
	if !on {
		t.Error("Expected led on, returned off")
	}

	bus.Bytes = append([]byte{1}, []byte("?T,19.5")...)
	v, err := e.GetTC()
	if err != nil {
		t.Error(err)
	}
	if v != 19.5 {
		t.Error("Expected 19.5 . Found:", v)
	}

	bus.Bytes = append([]byte{1}, []byte("?i,pH,2.8")...)
	d, i, err := e.Information()
	if err != nil {
		t.Error(err)
	}

	if d != "pH" {
		t.Error("Expected device pH. Found:", d)
	}

	if i != "2.8" {
		t.Error("Expected version 2.8. Found:", i)
	}

}

func TestEZOHalAdapter(t *testing.T) {
	factory := Factory()
	params := map[string]interface{}{}

	bus := i2c.MockBus()

	_, err := factory.NewDriver(params, bus)

	if err == nil {
		t.Error("EZO Driver creation should fail when configuration is invalid")
	}

	params["Address"] = 0x93

	e, err := factory.NewDriver(params, bus)
	if err != nil {
		t.Error(err)
	}
	d, ok := e.(hal.AnalogInputDriver)
	if !ok {
		t.Error("Failed to type cast ezo driver to ADC driver")
	}
	if d.Metadata().Name != _ezoName {
		t.Error("Unexpected name")
	}
	if !d.Metadata().HasCapability(hal.AnalogInput) {
		t.Error("PH Capability should exist")
	}
	if d.Metadata().HasCapability(hal.DigitalInput) {
		t.Error("Input Capability should not exist")
	}

	if len(d.AnalogInputPins()) != 1 {
		t.Error("Expected only one channel")
	}
	if _, err := d.AnalogInputPin(1); err == nil {
		t.Error("Expected error for invalid channel name")
	}

	if _, err := d.AnalogInputPin(0); err != nil {
		t.Error(err)
	}
	if err := d.Close(); err != nil {
		t.Error(err)
	}
}
