package ph_board

import (
	"fmt"
	"testing"

	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
)

func demo() {
	d, _ := NewDriver([]byte(`{"address":43}`), i2c.MockBus())
	ch, _ := d.AnalogInputPin(0)
	v, _ := ch.Read()
	fmt.Println(v)
}

func TestPhBoardDriver(t *testing.T) {
	bus := i2c.MockBus()
	bus.Bytes = make([]byte, 2)
	_, err := HalAdapter([]byte(""), bus)
	if err == nil {
		t.Error("Adapter creation should fail when json config is invalid")
	}
	configJSON := `
	{
		"address":16
	}
	`
	d, err := NewDriver([]byte(configJSON), bus)
	if err != nil {
		t.Error(err)
	}
	if d.Metadata().Name != "ph-board" {
		t.Error("Unexpected name")
	}
	if !d.Metadata().HasCapability(hal.AnalogInput) {
		t.Error("Analog input Capability should exist")
	}
	if d.Metadata().HasCapability(hal.DigitalInput) {
		t.Error("Digital Input Capability should not exist")
	}

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
	v, err := ch.Read()
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
