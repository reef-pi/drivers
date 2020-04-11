package sht3x

import (
	"errors"
	"fmt"
	"sync"

	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
)

const addressParam = "Address"

type factory struct {
	meta       hal.Metadata
	parameters []hal.ConfigParameter
}

var f *factory
var once sync.Once

// Factory returns a singleton pH board Driver factory
func Factory() hal.DriverFactory {
	once.Do(func() {
		f = &factory{
			meta: hal.Metadata{
				Name:         "sht31d",
				Description:  "SHT31D humidity and temperature sensor",
				Capabilities: []hal.Capability{hal.AnalogInput},
			},
			parameters: []hal.ConfigParameter{
				{
					Name:    addressParam,
					Type:    hal.Integer,
					Order:   0,
					Default: 0x44,
				},
			},
		}
	})
	return f
}

func (f *factory) Metadata() hal.Metadata {
	return f.meta
}

func (f *factory) GetParameters() []hal.ConfigParameter {
	return f.parameters
}

func (f *factory) ValidateParameters(parameters map[string]interface{}) (bool, map[string][]string) {
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

func (f *factory) NewDriver(parameters map[string]interface{}, hardwareResources interface{}) (hal.Driver, error) {
	if valid, failures := f.ValidateParameters(parameters); !valid {
		return nil, errors.New(hal.ToErrorString(failures))
	}
	intAddress, _ := hal.ConvertToInt(parameters[addressParam])
	address := byte(intAddress)
	bus := hardwareResources.(i2c.Bus)
	return NewDriver(address, bus, f.meta)
}
