package drivers

import (
	"testing"

	"github.com/reef-pi/rpi/i2c"
)

func TestHT16K33(t *testing.T) {
	bus := i2c.MockBus()
	bus.Bytes = make([]byte, 16, 16)
	h := NewHT16K33(bus)
	if err := h.Setup(); err != nil {
		t.Fatal(err)
	}
	if err := h.Display("REEF"); err != nil {
		t.Fatal(err)
	}
}
