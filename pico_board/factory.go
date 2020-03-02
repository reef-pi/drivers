package pico_board

import (
	"errors"
	"fmt"
	"sync"

	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
)

const addressParam = "Address"

type phFactory struct {
	meta       hal.Metadata
	parameters []hal.ConfigParameter
}

var factory *phFactory
var once sync.Once

// Factory returns a singleton pico-board Driver factory
func Factory() hal.DriverFactory {

	once.Do(func() {
		factory = &phFactory{
			meta: hal.Metadata{
				Name:         "pico-board",
				Description:  "Isolated ATSAMD10 pH driver on the blueAcro Pico board",
				Capabilities: []hal.Capability{hal.AnalogInput},
			},
			parameters: []hal.ConfigParameter{
				{
					Name:    addressParam,
					Type:    hal.Integer,
					Order:   0,
					Default: 0x45,
				},
			},
		}
	})

	return factory
}

func (f *phFactory) Metadata() hal.Metadata {
	return f.meta
}

func (f *phFactory) GetParameters() []hal.ConfigParameter {
	return f.parameters
}

func (f *phFactory) ValidateParameters(parameters map[string]interface{}) (bool, map[string][]string) {

	var failures = make(map[string][]string)

	if v, ok := parameters[addressParam]; ok {
		val, ok := hal.ConvertToInt(v)
		if !ok {
			failure := fmt.Sprint(addressParam, " is not a number. ", v, " was received.")
			failures[addressParam] = append(failures[addressParam], failure)
		}
		if val <= 0 || val >= 256 {
			failure := fmt.Sprint(addressParam, " is out of range (1 - 255). ", v, " was received.")
			failures[addressParam] = append(failures[addressParam], failure)
		}
	} else {
		failure := fmt.Sprint(addressParam, " is a required parameter, but was not received.")
		failures[addressParam] = append(failures[addressParam], failure)
	}

	return len(failures) == 0, failures
}

func (f *phFactory) NewDriver(parameters map[string]interface{}, hardwareResources interface{}) (hal.Driver, error) {
	if valid, failures := f.ValidateParameters(parameters); !valid {
		return nil, errors.New(hal.ToErrorString(failures))
	}

	address := byte(parameters[addressParam].(int))

	bus := hardwareResources.(i2c.Bus)

	if err := bus.WriteBytes(address, []byte{0x06}); err != nil {
		return nil, err
	}
	if err := bus.WriteBytes(address, []byte{0x40, 0x06}); err != nil {
		return nil, err
	}
	if err := bus.WriteBytes(address, []byte{0x08}); err != nil {
		return nil, err
	}

	ch, err := newChannel(bus, address)
	if err != nil {
		return nil, err
	}

	return &driver{
		channels: []hal.AnalogInputPin{ch},
		meta:     f.meta,
	}, nil
}
