package file

import (
	"errors"
	"fmt"
	"sync"

	"github.com/reef-pi/hal"
)

const pathParam = "Path"

type factory struct {
	meta       hal.Metadata
	parameters []hal.ConfigParameter
}

var analogFactory *factory
var once sync.Once

// AnalogFactory returns a singleton File based Analog Driver factory
func AnalogFactory() hal.DriverFactory {

	once.Do(func() {
		analogFactory = &factory{
			meta: hal.Metadata{
				Name:         "analog-file",
				Description:  "A simple file based analog hal driver",
				Capabilities: []hal.Capability{hal.AnalogInput},
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

	return analogFactory
}

func (f *factory) Metadata() hal.Metadata {
	return f.meta
}

func (f *factory) GetParameters() []hal.ConfigParameter {
	return f.parameters
}

func (f *factory) ValidateParameters(parameters map[string]interface{}) (bool, map[string][]string) {
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

func (f *factory) NewDriver(parameters map[string]interface{}, hardwareResources interface{}) (hal.Driver, error) {
	if valid, failures := f.ValidateParameters(parameters); !valid {
		return nil, errors.New(hal.ToErrorString(failures))
	}

	c, err := hal.CalibratorFactory([]hal.Measurement{})
	if err != nil {
		return nil, err
	}

	driver := &analog{
		path:       parameters[pathParam].(string),
		calibrator: c,
		meta:       f.meta,
	}
	return driver, nil
}
