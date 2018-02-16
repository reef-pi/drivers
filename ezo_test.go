package drivers

import (
	"github.com/reef-pi/rpi/i2c"
	"testing"
)

func TestEZO(t *testing.T) {
	bus := i2c.MockBus()
	e := NewAtlasEZO(byte(0x93), bus)
	if err := e.Read(); err != nil {
		t.Error(err)
	}
}
