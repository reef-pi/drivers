package tplink

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/reef-pi/hal"
)

const addressParam = "Address"

type HS103Plug struct {
	state   bool
	command *cmd
	meta    hal.Metadata
}

func newHS103Plug(addr string, meta hal.Metadata) *HS103Plug {
	return &HS103Plug{
		meta: meta,
		command: &cmd{
			addr: addr,
			cf: func(proto, addr string, t time.Duration) (Conn, error) {
				return net.DialTimeout(proto, addr, t)
			},
		},
	}
}

func (p *HS103Plug) SetFactory(cf ConnectionFactory) {
	p.command.cf = cf
}
func (p *HS103Plug) On() error {
	cmd := new(CmdRelayState)
	cmd.System.RelayState.State = 1
	if _, err := p.command.Execute(cmd, false); err != nil {
		return err
	}
	p.state = true
	return nil
}

func (p *HS103Plug) Off() error {
	cmd := new(CmdRelayState)
	cmd.System.RelayState.State = 0
	if _, err := p.command.Execute(cmd, false); err != nil {
		return err
	}
	p.state = false
	return nil
}

func (p *HS103Plug) Info() (*Sysinfo, error) {
	buf, err := p.command.Execute(new(Plug), true)
	if err != nil {
		return nil, err
	}
	var d Plug
	if err := json.Unmarshal(buf, &d); err != nil {
		return nil, err
	}
	return &d.System.Sysinfo, nil
}

func (p *HS103Plug) Metadata() hal.Metadata {
	return p.meta
}

func (p *HS103Plug) Name() string {
	return p.meta.Name
}

func (p *HS103Plug) Number() int {
	return 0
}
func (p *HS103Plug) DigitalOutputPins() []hal.DigitalOutputPin {
	return []hal.DigitalOutputPin{p}
}

func (p *HS103Plug) DigitalOutputPin(i int) (hal.DigitalOutputPin, error) {
	if i != 0 {
		return nil, fmt.Errorf("invalid pin: %d", i)
	}
	return p, nil
}

func (p *HS103Plug) Write(state bool) error {
	if state {
		return p.On()
	}
	return p.Off()
}

func (p *HS103Plug) LastState() bool {
	return p.state
}

func (p *HS103Plug) Close() error {
	return nil
}
func (p *HS103Plug) Pins(cap hal.Capability) ([]hal.Pin, error) {
	switch cap {
	case hal.DigitalOutput:
		return []hal.Pin{p}, nil
	default:
		return nil, fmt.Errorf("unsupported capability:%s", cap.String())
	}
}

type hs103Factory struct {
	meta       hal.Metadata
	parameters []hal.ConfigParameter
}

var factory103 *hs103Factory
var hs103once sync.Once

// HS103Factory returns a singleton HS103 Driver factory
func HS103Factory() hal.DriverFactory {

	hs103once.Do(func() {
		factory103 = &hs103Factory{
			meta: hal.Metadata{
				Name:        "tplink-hs103",
				Description: "tplink hs103 series smart plug driver",
				Capabilities: []hal.Capability{
					hal.DigitalOutput,
				},
			},
			parameters: []hal.ConfigParameter{
				{
					Name:    addressParam,
					Type:    hal.String,
					Order:   0,
					Default: "192.168.1.11:9999",
				},
			},
		}
	})

	return factory103
}

func (f *hs103Factory) Metadata() hal.Metadata {
	return f.meta
}

func (f *hs103Factory) GetParameters() []hal.ConfigParameter {
	return f.parameters
}

func (f *hs103Factory) ValidateParameters(parameters map[string]interface{}) (bool, map[string][]string) {

	var failures = make(map[string][]string)

	if v, ok := parameters[addressParam]; ok {
		_, ok := v.(string)
		if !ok {
			failure := fmt.Sprint(addressParam, " is not a string. ", v, " was received.")
			failures[addressParam] = append(failures[addressParam], failure)
		}
	} else {
		failure := fmt.Sprint(addressParam, " is a required parameter, but was not received.")
		failures[addressParam] = append(failures[addressParam], failure)
	}

	return len(failures) == 0, failures
}

func (f *hs103Factory) NewDriver(parameters map[string]interface{}, _ interface{}) (hal.Driver, error) {
	if valid, failures := f.ValidateParameters(parameters); !valid {
		return nil, errors.New(hal.ToErrorString(failures))
	}

	addr := parameters[addressParam].(string)

	return newHS103Plug(addr, f.meta), nil
}
