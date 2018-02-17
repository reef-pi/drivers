package drivers

import (
	"github.com/reef-pi/rpi/i2c"
	"testing"
)

func TestEZO(t *testing.T) {
	bus := i2c.MockBus()
	e := NewAtlasEZO(byte(0x93), bus)
	bus.Bytes = []byte("19.65")
	if _, err := e.Read(); err != nil {
		t.Error(err)
	}
	bus.Bytes = []byte("09.65")
	if _, err := e.Read(); err == nil {
		t.Error("Values starting with 0 should fail")
	}
}
