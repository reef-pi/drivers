package esp32

import (
	"bytes"
	"github.com/reef-pi/hal"
	"io"
	"net/http"
	"testing"
)

func NopClient(r *http.Request) (*http.Response, error) {
	resp := new(http.Response)
	resp.Body = io.NopCloser(bytes.NewBuffer([]byte("2")))
	resp.StatusCode = 200
	return resp, nil
}

func TestESP32Driver(t *testing.T) {
	f := FactoryWithClient(NopClient)
	params := map[string]interface{}{
		"Address":        "192.168.86.2",
		"digital-output": "2,3",
		"digital-input":  "1,4",
		"analog-input":   "6,7",
		"pwm":            "10",
	}

	d, err := f.NewDriver(params, nil)
	if err != nil {
		t.Error(err)
	}

	do, ok := d.(hal.DigitalOutputDriver)
	if !ok {
		t.Error("failed to typecast to digital output driver")
		return
	}
	oPins := do.DigitalOutputPins()

	if len(oPins) != 2 {
		t.Error("expected 2 digital output pins, found:", len(oPins))
		return
	}
	if err := oPins[0].Write(true); err != nil {
		t.Error(err)
	}
	di, ok := d.(hal.DigitalInputDriver)
	if !ok {
		t.Error("failed to typecast to digital input driver")
		return
	}
	iPins := di.DigitalInputPins()

	if len(iPins) != 2 {
		t.Error("expected 2 digital input pins, found:", len(iPins))
		return
	}
	if _, err := iPins[0].Read(); err != nil {
		t.Error(err)
	}

	pd, ok := d.(hal.PWMDriver)
	if !ok {
		t.Error("failed to typecast to pwm driver")
		return
	}
	chs := pd.PWMChannels()
	if len(chs) != 1 {
		t.Error("expected 1 pwm pin found:", len(iPins))
		return
	}
	if err := chs[0].Set(50); err != nil {
		t.Error(err)
	}

	ad, ok := d.(hal.AnalogInputDriver)
	if !ok {
		t.Error("failed to typecast to analog input driver")
		return
	}
	aPins := ad.AnalogInputPins()
	if len(aPins) != 2 {
		t.Error("expected 2 analog input pin found:", len(aPins))
		return
	}
	if _, err := aPins[0].Value(); err != nil {
		t.Error(err)
	}
}
