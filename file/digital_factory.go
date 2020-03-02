package file

import (
	"errors"
	"fmt"
	"sync"

	"github.com/reef-pi/hal"
)

type dFactory struct {
	meta       hal.Metadata
	parameters []hal.ConfigParameter
}

var digitalFactory *dFactory
var digitalOnce sync.Once

// DigitalFactory returns a singleton File based Digital Driver factory
func DigitalFactory() hal.DriverFactory {

	digitalOnce.Do(func() {
		digitalFactory = &dFactory{
			meta: hal.Metadata{
				Name:         "digital-file",
				Description:  "A simple file based digital hal driver",
				Capabilities: []hal.Capability{hal.DigitalInput, hal.DigitalOutput, hal.PWM},
			},
			parameters: []hal.ConfigParameter{
				{
					Name:    pathParam,
					Type:    hal.String,
					Order:   0,
					Default: "/path/to/file",
				},
			},
		}
	})

	return digitalFactory
}

func (f *dFactory) Metadata() hal.Metadata {
	return f.meta
}

func (f *dFactory) GetParameters() []hal.ConfigParameter {
	return f.parameters
}

func (f *dFactory) ValidateParameters(parameters map[string]interface{}) (bool, map[string][]string) {

	var failures = make(map[string][]string)

	if v, ok := parameters[pathParam]; ok {
		val, ok := v.(string)
		if !ok {
			failure := fmt.Sprint(pathParam, " is not a string. ", v, " was received.")
			failures[pathParam] = append(failures[pathParam], failure)
		}
		if len(val) < 1 {
			failure := fmt.Sprint(pathParam, " not long enough to be valid. ", v, " was received.")
			failures[pathParam] = append(failures[pathParam], failure)
		}
	} else {
		failure := fmt.Sprint(pathParam, " is required parameter, but was not received.")
		failures[pathParam] = append(failures[pathParam], failure)
	}

	return len(failures) == 0, failures
}

func (f *dFactory) NewDriver(parameters map[string]interface{}, hardwareResources interface{}) (hal.Driver, error) {
	if valid, failures := f.ValidateParameters(parameters); !valid {
		return nil, errors.New(hal.ToErrorString(failures))
	}

	driver := &digital{
		path: parameters[pathParam].(string),
		meta: f.meta,
	}

	return driver, nil
}
