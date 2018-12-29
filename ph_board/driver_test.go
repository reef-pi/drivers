package ph_board

import (
	"fmt"
	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
	"testing"
)

func demo() {
	d, _ := HalAdapter([]byte(`{"address":43}`), i2c.MockBus())
	ch, _ := d.ADCChannel("0")
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
	d, err := HalAdapter([]byte(configJSON), bus)
	if err != nil {
		t.Error(err)
	}
	if d.Metadata().Name != "ph-board" {
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
	if _, err := d.ADCChannel("1"); err == nil {
		t.Error("Expected error for invalid channel name")
	}

	ch, err := d.ADCChannel("0")
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
