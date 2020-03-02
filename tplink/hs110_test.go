package tplink

import (
	"testing"

	"github.com/reef-pi/hal"
)

func TestHS110Plug(t *testing.T) {
	p := newHS110Plug("127.0.0.1:9999", hal.Metadata{})
	nop := NewNop()
	p.SetFactory(nop.Factory)
	if err := p.On(); err != nil {
		t.Error(err)
	}
	nop.Buffer([]byte(`{}`))
	if err := p.Off(); err != nil {
		t.Error(err)
	}

	f := HS110Factory()

	params := map[string]interface{}{
		"Address": "http://192.168.1.5:3000",
	}

	d, err := f.NewDriver(params, nil)

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
