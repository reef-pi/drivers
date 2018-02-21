package drivers

import (
	"github.com/reef-pi/rpi/i2c"
	"testing"
)

func TestSHT31D(t *testing.T) {
	bus := i2c.MockBus()
	s := NewSHT31D(byte(0x93), bus)
	s.delay = 0
	bus.Bytes = make([]byte, 6)
	if _, _, err := s.Read(); err == nil {
		t.Error("Expect crc failure")
	}
	if err := s.SetAccuracy(SHT3xAccuracyLow); err != nil {
		t.Error(err)
	}
	if _, err := s.SerialNumber(); err == nil {
		t.Error("Expect crc failure")
	}
}
