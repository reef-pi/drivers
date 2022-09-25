package esp32

import (
	"errors"
	"fmt"
	"github.com/reef-pi/hal"
	"strconv"
	"strings"
	"sync"
)

type factory struct {
	meta       hal.Metadata
	parameters []hal.ConfigParameter
}

var esp32DriverFactory *factory
var once sync.Once

const Address = "Address"
const PWMPins = "PWM Pins"

func Factory() hal.DriverFactory {

	once.Do(func() {
		esp32DriverFactory = &factory{
			meta: hal.Metadata{
				Name:        "reef-pi ESP32 driver",
				Description: "Simple HTTP based full featured HAL driver for reef-pi",
				Capabilities: []hal.Capability{
					hal.PWM,
					hal.DigitalOutput,
					hal.DigitalInput,
					hal.AnalogInput,
				},
			},
			parameters: []hal.ConfigParameter{
				{
					Name:    Address,
					Type:    hal.String,
					Order:   0,
					Default: "192.1.168.4",
				},
				{
					Name:    PWMPins,
					Type:    hal.String,
					Order:   1,
					Default: "",
				},
				{
					Name:    hal.DigitalOutput.String(),
					Type:    hal.String,
					Order:   2,
					Default: "",
				},
				{
					Name:    hal.DigitalInput.String(),
					Type:    hal.String,
					Order:   3,
					Default: "",
				},
				{
					Name:    hal.AnalogInput.String(),
					Type:    hal.String,
					Order:   4,
					Default: "",
				},
			},
		}
	})

	return esp32DriverFactory
}

func (f *factory) GetParameters() []hal.ConfigParameter {
	return f.parameters
}

func (f *factory) ValidateParameters(parameters map[string]interface{}) (bool, map[string][]string) {
	var failures = make(map[string][]string)
	pins := make(map[int]struct{})

	if v, ok := parameters[Address]; ok {
		val, ok := v.(string)
		if !ok {
			failure := fmt.Sprint(Address, " is not a string. ", v, " was received.")
			failures[Address] = append(failures[Address], failure)
		} else if len(val) <= 0 {
			failure := fmt.Sprint(Address, " empty values are not allowed.")
			failures[Address] = append(failures[Address], failure)
		} else if len(val) >= 256 {
			failure := fmt.Sprint(Address, " size should be lower than 255 characters. ", val, " was received.")
			failures[Address] = append(failures[Address], failure)
		}
	} else {
		failure := fmt.Sprint(Address, " is a required parameter, but was not received.")
		failures[Address] = append(failures[Address], failure)
	}
	for _, cap := range []hal.Capability{hal.DigitalOutput, hal.DigitalInput, hal.PWM, hal.AnalogInput} {
		if v, ok := parameters[cap.String()]; ok {
			val, ok := v.(string)
			if !ok {
				failure := fmt.Sprint(cap, " is not a string. ", parameters[cap.String()], " was received.")
				failures[cap.String()] = append(failures[cap.String()], failure)
			}
			if val != "" {
				sPins := strings.Split(val, ",")
				for _, s := range sPins {
					i, err := strconv.Atoi(s)
					if err != nil {
						failures[cap.String()] = append(failures[cap.String()], fmt.Sprint(cap, " pin", s, " is not an integer"))
					}
					_, ok := pins[i]
					if ok {
						failures[cap.String()] = append(failures[cap.String()], fmt.Sprint(cap, " pin", s, " is already in use"))
					}
					pins[i] = struct{}{}
				}
			}
		}
	}
	return len(failures) == 0, failures
}

func (f *factory) Metadata() hal.Metadata {
	return f.meta
}

func (f *factory) NewDriver(parameters map[string]interface{}, hardwareResources interface{}) (hal.Driver, error) {
	if valid, failures := f.ValidateParameters(parameters); !valid {
		return nil, errors.New(hal.ToErrorString(failures))
	}
	pins := make(map[hal.Capability][]int)
	for _, cap := range []hal.Capability{hal.DigitalOutput, hal.DigitalInput, hal.PWM, hal.AnalogInput} {
		if v, ok := parameters[cap.String()]; ok {
			val, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("failed to type cast '%s' parameter value '%v' as string", cap, v)
			}
			if val != "" {
				sPins := strings.Split(val, ",")
				for _, s := range sPins {
					i, err := strconv.Atoi(s)
					if err != nil {
						return nil, fmt.Errorf("failed to convert '%s' pin '%v' to integrer. Error:%w", cap, s, err)
					}
					pins[cap] = append(pins[cap], i)
				}
			}
		}
	}
	return &driver{
		meta:    f.meta,
		address: parameters[Address].(string),
		pins:    pins,
	}, nil
}
