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
	d, err := NewAnalog(temp.Name())
	if err != nil {
		t.Fatal(err)
	}
	meta := d.Metadata()
	if len(meta.Capabilities) != 1 {
		t.Error("Expected 1 capabilities, found:", len(meta.Capabilities))
	}
	dig := hal.ADCDriver(d)
	if len(dig.ADCChannels()) != 1 {
		t.Error("Expected a single input pin, found:", len(dig.ADCChannels()))
	}
	pin, err := dig.ADCChannel(0)
	if err != nil {
		t.Error(err)
	}
	if _, err := pin.Read(); err != nil {
		t.Error(err)
	}
}
