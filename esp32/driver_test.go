package esp32

import (
	"github.com/reef-pi/hal"
	"testing"
)

func TestESP32Driver(t *testing.T) {
	f := Factory()
	params := map[string]interface{}{
		"Address":        "192.168.86.2",
		"digital-output": "2,3",
	}

	driver, err := f.NewDriver(params, nil)
	if err != nil {
		t.Error(err)
	}
	d, ok := driver.(hal.DigitalInputDriver)
	if !ok {
		t.Error("failed to typecast to digital driver")
		return
	}
	pins := d.DigitalInputPins()

	if len(pins) != 2 {
		t.Error("expected 2 digital pins, found:", len(pins))
		return
	}
	if _, err := pins[0].Read(); err != nil {
		t.Error(err)
	}
}
