package drivers

import (
	"github.com/reef-pi/rpi/i2c"
	"testing"
)

func TestPCA9685(t *testing.T) {
	bus := i2c.MockBus()
	p := NewPCA9685(0x70, bus)
	if err := p.Wake(); err != nil {
		t.Fatal(err)
	}
	p.SetPwm(10, 0, 10)
}
