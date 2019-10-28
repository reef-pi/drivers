package tplink

import (
	"testing"

	"github.com/reef-pi/hal"
)

func TestHS110Plug(t *testing.T) {
	p := NewHS110Plug("127.0.0.1:9999")
	nop := NewNop()
	p.cnFactory = nop.Factory
	if err := p.On(); err != nil {
		t.Error(err)
	}
	nop.Buffer([]byte(`{}`))
	if err := p.Off(); err != nil {
		t.Error(err)
	}

	d, err := HS110HALAdapter([]byte(`{"address":"127.0.0.1:3000"}`), nil)
	if err != nil {
		t.Error(err)
	}
	if d.Metadata().Name == "" {
		t.Error("HAL metadata should not have empty name")
	}

	d1, ok := d.(hal.DigitalOutputDriver)
	if !ok {
		t.Fatal("Failed to type cast to output driver")
	}

	if len(d1.DigitalOutputPins()) != 1 {
		t.Error("Expected exactly one output pin")
	}
	pin, err := d1.DigitalOutputPin(0)
	if err != nil {
		t.Error(err)
	}
	if pin.LastState() != false {
		t.Error("Expected initial state to be false")
	}
}
