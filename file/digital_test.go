package file

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/reef-pi/hal"
)

func TestDigitalInput(t *testing.T) {
	temp, err := ioutil.TempFile("", "hal-file-driver-testing")
	if err != nil {
		t.Fatal(err)
	}
	temp.Close()
	defer os.Remove(temp.Name())

	params := map[string]interface{}{
		"Path": temp.Name(),
	}

	f := DigitalFactory()
	d, err := f.NewDriver(params, nil)
	if err != nil {
		t.Error(err)
	}

	meta := d.Metadata()
	if len(meta.Capabilities) != 3 {
		t.Error("Expected 3 capabilities, found:", len(meta.Capabilities))
	}

	dig, ok := d.(hal.DigitalInputDriver)
	if !ok {
		t.Error("Failed to type cast digital file driver to digital input driver")
	}

	if len(dig.DigitalInputPins()) != 1 {
		t.Error("Expected a single input pin, found:", len(dig.DigitalInputPins()))
	}
	pin, err := dig.DigitalInputPin(0)
	if err != nil {
		t.Error(err)
	}
	b, err := pin.Read()
	if err != nil {
		t.Error(err)
	}
	if b {
		t.Error("Expected false , found true")
	}
}
func TestDigitalOutput(t *testing.T) {
	temp, err := ioutil.TempFile("", "hal-file-driver-testing")
	if err != nil {
		t.Fatal(err)
	}
	temp.Close()
	defer os.Remove(temp.Name())

	params := map[string]interface{}{
		"Path": temp.Name(),
	}

	f := DigitalFactory()
	d, err := f.NewDriver(params, nil)
	if err != nil {
		t.Error(err)
	}

	dig, ok := d.(hal.DigitalOutputDriver)
	if !ok {
		t.Error("Failed to type cast digital file driver to digital output driver")
	}

	pin, err := dig.DigitalOutputPin(0)
	if err != nil {
		t.Error(err)
	}
	if err := pin.Write(true); err != nil {
		t.Error(err)
	}
}

func TestPWMOutput(t *testing.T) {
	temp, err := ioutil.TempFile("", "hal-file-driver-testing")
	if err != nil {
		t.Fatal(err)
	}
	temp.Close()
	defer os.Remove(temp.Name())

	params := map[string]interface{}{
		"Path": temp.Name(),
	}

	f := DigitalFactory()
	d, err := f.NewDriver(params, nil)
	if err != nil {
		t.Error(err)
	}

	dig, ok := d.(hal.PWMDriver)
	if !ok {
		t.Error("Failed to type cast digital file driver to digital output driver")
	}

	pin, err := dig.PWMChannel(0)
	if err != nil {
		t.Error(err)
	}
	if err := pin.Set(12.3); err != nil {
		t.Error(err)
	}
}
