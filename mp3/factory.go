package mp3

import (
	"errors"
	"fmt"
	"github.com/reef-pi/hal"
	"sync"
)

const (
	fileParam = "File"
	loopParam = "Loop"
)

type factory struct {
	meta       hal.Metadata
	parameters []hal.ConfigParameter
}

var f *factory
var once sync.Once

// Factory returns a singleton mp3 driver factory
func Factory() hal.DriverFactory {
	once.Do(func() {
		f = &factory{
			meta: hal.Metadata{
				Name:        "mp3",
				Description: "MP3 Audio as digital output driver",
				Capabilities: []hal.Capability{
					hal.DigitalOutput,
				},
			},
			parameters: []hal.ConfigParameter{
				{
					Name:    fileParam,
					Type:    hal.String,
					Order:   0,
					Default: "/var/lib/reef-pi/mp3/alert.mp3",
				},
				{
					Name:    loopParam,
					Type:    hal.Boolean,
					Order:   1,
					Default: false,
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

	if v, ok := parameters[fileParam]; ok {
		_, ok := v.(string)
		if !ok {
			failure := fmt.Sprint(fileParam, " is not a string. ", v, " was received.")
			failures[fileParam] = append(failures[fileParam], failure)
		}
	} else {
		failure := fmt.Sprint(fileParam, " is a required parameter, but was not received.")
		failures[fileParam] = append(failures[fileParam], failure)
	}

	if v, ok := parameters[loopParam]; ok {
		_, ok := v.(bool)
		if !ok {
			failure := fmt.Sprint(loopParam, " is not a bool. ", v, " was received.")
			failures[loopParam] = append(failures[loopParam], failure)
		}
	} else {
		failure := fmt.Sprint(loopParam, " is a required parameter, but was not received.")
		failures[loopParam] = append(failures[loopParam], failure)
	}

	return len(failures) == 0, failures
}

func (f *factory) NewDriver(parameters map[string]interface{}, _ interface{}) (hal.Driver, error) {
	if valid, failures := f.ValidateParameters(parameters); !valid {
		return nil, errors.New(hal.ToErrorString(failures))
	}

	file := parameters[fileParam].(string)
	loop := parameters[loopParam].(bool)

	return &Driver{
		meta: f.meta,
		conf: Config{
			File: file,
			Loop: loop,
		},
	}, nil
}
