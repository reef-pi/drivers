package pca9685

import (
	"testing"

	"github.com/reef-pi/rpi/i2c"
)

func TestNew(t *testing.T) {
	bus := i2c.MockBus()
	p := New(0x70, bus)
	if err := p.Wake(); err != nil {
		t.Fatal(err)
	}
	p.SetPwm(10, 0, 10)
}
