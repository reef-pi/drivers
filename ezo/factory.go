package ezo

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
)

type factory struct {
	meta       hal.Metadata
	parameters []hal.ConfigParameter
}

var ezoFactory *factory
var once sync.Once

// Factory returns a singleton EZO Driver factory
func Factory() hal.DriverFactory {

	once.Do(func() {
		ezoFactory = &factory{
			meta: hal.Metadata{
				Name:         _ezoName,
				Description:  "Atlas Scientific EZO board for pH sensor",
				Capabilities: []hal.Capability{hal.AnalogInput},
			},
			parameters: []hal.ConfigParameter{
				{
					Name:    "Address",
					Type:    hal.Integer,
					Order:   0,
					Default: 68,
				},
			},
		}
	})

	return ezoFactory
}

func (f *factory) Metadata() hal.Metadata {
	return f.meta
}

//Implement hal.Driver interface
func (f *factory) GetParameters() []hal.ConfigParameter {
	return f.parameters
}

func (f *factory) ValidateParameters(parameters map[string]interface{}) (bool, map[string][]string) {

	var failures = make(map[string][]string)

	if address, ok := parameters["Address"]; ok {
		val, ok := hal.ConvertToInt(address)
		if !ok {
			failure := fmt.Sprint("Address is not an integer. ", address, " was received.")
			failures["Address"] = append(failures["Sample Parameter"], failure)
		}
		if val < 0 || val > 255 {
			failure := fmt.Sprint("Address is out of range. It should be between 0 and 255, but ", address, " was received.")
			failures["Address"] = append(failures["Sample Parameter"], failure)
		}
	} else {
		failure := fmt.Sprint("Address is not a required parameter, but was not found.")
		failures["Address"] = append(failures["Sample Parameter"], failure)
	}

	return len(failures) == 0, failures
}

func (f *factory) NewDriver(parameters map[string]interface{}, hardwareResources interface{}) (hal.Driver, error) {
	if valid, failures := f.ValidateParameters(parameters); !valid {
		return nil, errors.New(hal.ToErrorString(failures))
	}

	address, _ := hal.ConvertToInt(parameters["Address"])

	driver := &AtlasEZO{
		addr:  byte(address),
		bus:   hardwareResources.(i2c.Bus),
		delay: time.Second,
		meta: hal.Metadata{
			Name:         _ezoName,
			Description:  "Atlas Scientific EZO board for pH sensor",
			Capabilities: []hal.Capability{hal.AnalogInput},
		},
	}

	return driver, nil
}
