package file

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/reef-pi/hal"
)

func TestAnalogInput(t *testing.T) {
	temp, err := ioutil.TempFile("", "hal-file-driver-testing")
	if err != nil {
		t.Fatal(err)
	}
	temp.Write([]byte("23.1"))
	temp.Close()
	defer os.Remove(temp.Name())

	params := map[string]interface{}{
		"Path": temp.Name(),
	}

	f := AnalogFactory()
	d, err := f.NewDriver(params, nil)

	if err != nil {
		t.Fatal(err)
	}
	meta := d.Metadata()
	if len(meta.Capabilities) != 1 {
		t.Error("Expected 1 capabilities, found:", len(meta.Capabilities))
	}

	dig, ok := d.(hal.AnalogInputDriver)
	if !ok {
		t.Error("Failed to type cast analog file driver to analog input driver")
	}

	if len(dig.AnalogInputPins()) != 1 {
		t.Error("Expected a single input pin, found:", len(dig.AnalogInputPins()))
	}
	pin, err := dig.AnalogInputPin(0)
	if err != nil {
		t.Error(err)
	}
	if _, err := pin.Value(); err != nil {
		t.Error(err)
	}
}
