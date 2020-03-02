package tplink

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/reef-pi/hal"
)

type (
	EmeterCmd struct {
		Emeter struct {
			Realtime Realtime `json:"get_realtime"`
		} `json:"emeter"`
	}
	Realtime struct {
		Current  float64 `json:"current,omitempty"`
		Voltage  float64 `json:"voltage,omitempty"`
		Power    float64 `json:"power,omitempty"`
		Total    float64 `json:"total,omitempty"`
		ErrrCode int     `json:"err_code,omitempty"`
	}
)

type HS110Plug struct {
	HS103Plug
	calibrator hal.Calibrator
}

func (p *HS110Plug) Number() int {
	return 0
}

func newHS110Plug(addr string, meta hal.Metadata) *HS110Plug {
	cal, _ := hal.CalibratorFactory([]hal.Measurement{})

	return &HS110Plug{
		HS103Plug: HS103Plug{
			command: &cmd{
				addr: addr,
				cf:   TCPConnFactory,
			},
			meta: meta,
		},
		calibrator: cal,
	}
}

func (p *HS110Plug) RTEmeter() (*Realtime, error) {
	d, err := p.command.Execute(new(EmeterCmd), true)
	if err != nil {
		return nil, err
	}
	var cmd EmeterCmd
	if err := json.Unmarshal(d, &cmd); err != nil {
		return nil, err
	}
	return &cmd.Emeter.Realtime, nil
}

func (p *HS110Plug) SetFactory(cf ConnectionFactory) {
	p.command.cf = cf
}

func (p *HS110Plug) AnalogInputPins() []hal.AnalogInputPin {
	return []hal.AnalogInputPin{p}
}

func (p *HS110Plug) AnalogInputPin(i int) (hal.AnalogInputPin, error) {
	if i != 0 {
		return nil, fmt.Errorf("invalid channel number: %d", i)
	}
	return p, nil
}

func (p *HS110Plug) Read() (float64, error) {
	em, err := p.RTEmeter()
	if err != nil {
		return 0, err
	}
	return em.Current, nil
}

func (p *HS110Plug) Calibrate(points []hal.Measurement) error {
	cal, err := hal.CalibratorFactory(points)
	if err != nil {
		return err
	}
	p.calibrator = cal
	return nil
}
func (p *HS110Plug) Measure() (float64, error) {
	v, err := p.Read()
	if err != nil {
		return 0, err
	}
	if p.calibrator == nil {
		return 0, fmt.Errorf("Not calibrated")
	}
	return p.calibrator.Calibrate(v), nil
}

func (p *HS110Plug) Pins(cap hal.Capability) ([]hal.Pin, error) {
	switch cap {
	case hal.DigitalOutput, hal.AnalogInput:
		return []hal.Pin{p}, nil
	default:
		return nil, fmt.Errorf("unsupported capability:%s", cap)
	}
}

type hs110Factory struct {
	meta       hal.Metadata
	parameters []hal.ConfigParameter
}

var factory110 *hs110Factory
var hs110once sync.Once

// HS110Factory returns a singleton HS110 Driver factory
func HS110Factory() hal.DriverFactory {

	hs110once.Do(func() {
		factory110 = &hs110Factory{
			meta: hal.Metadata{
				Name:        "tplink-hs110",
				Description: "tplink hs110 series smart plug driver with current monitoring",
				Capabilities: []hal.Capability{
					hal.DigitalOutput, hal.AnalogInput,
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

	return factory110
}

func (f *hs110Factory) Metadata() hal.Metadata {
	return f.meta
}

func (f *hs110Factory) GetParameters() []hal.ConfigParameter {
	return f.parameters
}

func (f *hs110Factory) ValidateParameters(parameters map[string]interface{}) (bool, map[string][]string) {

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

func (f *hs110Factory) NewDriver(parameters map[string]interface{}, hardwareResources interface{}) (hal.Driver, error) {
	if valid, failures := f.ValidateParameters(parameters); !valid {
		return nil, errors.New(hal.ToErrorString(failures))
	}

	addr := parameters[addressParam].(string)

	return newHS110Plug(addr, f.meta), nil
}
