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
	d := NewDigital(temp.Name())
	meta := d.Metadata()
	if len(meta.Capabilities) != 3 {
		t.Error("Expected 3 capabilities, found:", len(meta.Capabilities))
	}
	dig := hal.DigitalInputDriver(d)
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
	d := NewDigital(temp.Name())
	dig := hal.DigitalOutputDriver(d)
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
	d := NewDigital(temp.Name())
	dig := hal.PWMDriver(d)
	pin, err := dig.PWMChannel(0)
	if err != nil {
		t.Error(err)
	}
	if err := pin.Set(12.3); err != nil {
		t.Error(err)
	}
}
