package execstream

import (
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"

	"encoding/json"
	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
)

type digitalPin struct {
	pin       int
	stream    streamSupervisor
	meta      hal.Metadata
	lastState bool
}

type streamSupervisor struct {
	path string
	cmd exec.Cmd
	meta hal.Metadata
}

func HalDigitalAdapter(c []byte, _ i2c.Bus) (hal.Driver, error) {
	var config Config
	if err := json.Unmarshal(c, &config); err != nil {
		return nil, err
	}
	return NewDigital(config.Address), nil
}

func NewDigital(p string) *streamSupervisor {
	return &streamSupervisor{
		path: p,
		meta: hal.Metadata{
			Name:         "executable-stream",
			Description:  "Driver to execute programs for I/O",
			Capabilities: []hal.Capability{hal.Input, hal.Output, hal.PWM},
		},
	}
}

func (f *streamSupervisor) Metadata() hal.Metadata {
	return f.meta
}

func (f *streamSupervisor) Close() error {
	return nil
}

func (f *streamSupervisor) Name() string {
	return f.path
}

func (f *digitalPin) Read() (bool, error) {
	data, err := ioutil.ReadFile(f.path)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(data)) == "1", nil
}

func (f *digitalPin) LastState() bool {
	return f.lastState
}

func (f *digitalPin) Write(b bool) error {
	f.lastState = b
	if b {
		return ioutil.WriteFile(f.path, []byte("1"), 0644)
	}
	return ioutil.WriteFile(f.path, []byte("0"), 0644)

}
func (f *digitalPin) Set(v float64) error {
	return ioutil.WriteFile(f.path, []byte(strconv.FormatFloat(v, 'f', -1, 64)), 0644)
}

func (f *streamSupervisor) InputPins() []hal.InputPin {
	return []hal.InputPin{f}
}

func (f *streamSupervisor) InputPin(_ int) (hal.InputPin, error) {
	return f, nil
}

func (f *streamSupervisor) OutputPins() []hal.OutputPin {
	return []hal.OutputPin{f}
}

func (f *streamSupervisor) OutputPin(_ int) (hal.OutputPin, error) {
	return f, nil
}
func (f *streamSupervisor) PWMChannels() []hal.PWMChannel {
	return []hal.PWMChannel{f}
}

func (f *streamSupervisor) PWMChannel(_ int) (hal.PWMChannel, error) {
	return f, nil
}
