package drivers

import (
	"github.com/reef-pi/hal"
	"testing"

	"github.com/reef-pi/rpi/i2c"
)

func TestEZO(t *testing.T) {
	bus := i2c.MockBus()
	e := NewAtlasEZO(byte(0x93), bus)
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
	bus := i2c.MockBus()
	_, err := EzoHalAdapter([]byte(""), bus)
	if err == nil {
		t.Error("Adapter creation should fail when json config is invalid")
	}
	configJSON := `
	{
		"address":16
	}
	`
	e, err := EzoHalAdapter([]byte(configJSON), bus)
	if err != nil {
		t.Error(err)
	}
	d, ok := e.(hal.ADCDriver)
	if !ok {
		t.Error("Failed to type cast ezo driver to ADC driver")
	}
	if d.Metadata().Name != _ezoName {
		t.Error("Unexpected name")
	}
	if !d.Metadata().HasCapability(hal.PH) {
		t.Error("PH Capability should exist")
	}
	if d.Metadata().HasCapability(hal.Input) {
		t.Error("Input Capability should not exist")
	}

	if len(d.ADCChannels()) != 1 {
		t.Error("Expected only one channel")
	}
	if _, err := d.ADCChannel(1); err == nil {
		t.Error("Expected error for invalid channel name")
	}

	if _, err := d.ADCChannel(0); err != nil {
		t.Error(err)
	}
	if err := d.Close(); err != nil {
		t.Error(err)
	}
}
