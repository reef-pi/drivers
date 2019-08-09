package file

import (
	"errors"
	"fmt"
	"sync"

	"github.com/reef-pi/hal"
)

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
					Name:    "Path",
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

	if v, ok := parameters["Path"]; ok {
		val, ok := v.(string)
		if !ok {
			failure := fmt.Sprint("Path is not a string. ", v, " was received.")
			failures["Path"] = append(failures["Path"], failure)
		}
		if len(val) < 1 {
			failure := fmt.Sprint("File path not long enough to be valid. ", v, " was received.")
			failures["Path"] = append(failures["Path"], failure)
		}
	} else {
		failure := fmt.Sprint("Path is required parameter, but was not received.")
		failures["Path"] = append(failures["Path"], failure)
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
		path:       parameters["Path"].(string),
		calibrator: c,
		meta:       f.meta,
	}
	return driver, nil
}
