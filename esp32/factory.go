package esp32

import (
	"errors"
	"fmt"
	"github.com/reef-pi/hal"
	"net/http"
	"strings"
	"sync"
)

type factory struct {
	meta       hal.Metadata
	parameters []hal.ConfigParameter
	client     HTTPClient
}

var esp32DriverFactory *factory
var once sync.Once

const Address = "Address"

func cap2string(c hal.Capability) string {
	return strings.Title(c.String())
}

func Factory() hal.DriverFactory {
	client := http.Client{Timeout: _timeout}
	return FactoryWithClient(client.Do)
}
func FactoryWithClient(c HTTPClient) hal.DriverFactory {

	once.Do(func() {
		esp32DriverFactory = &factory{
			client: c,
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
					Name:    cap2string(hal.DigitalOutput),
					Type:    hal.Integer,
					Order:   1,
					Default: 6,
				},
				{
					Name:    cap2string(hal.DigitalInput),
					Type:    hal.Integer,
					Order:   2,
					Default: 4,
				},
				{
					Name:    cap2string(hal.PWM),
					Type:    hal.Integer,
					Order:   3,
					Default: 4,
				},
				{
					Name:    cap2string(hal.AnalogInput),
					Type:    hal.Integer,
					Order:   4,
					Default: 2,
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
	for _, c := range []hal.Capability{hal.DigitalOutput, hal.DigitalInput, hal.PWM, hal.AnalogInput} {
		if v, ok := parameters[cap2string(c)]; ok {
			val, converted := hal.ConvertToInt(v)
			if !converted {
				failure := fmt.Sprint(c, " is not an integer. ", parameters[cap2string(c)], " was received")
				failures[cap2string(c)] = append(failures[cap2string(c)], failure)
			}
			if val < 0 {
				failures[cap2string(c)] = append(failures[cap2string(c)], fmt.Sprint(c, " pin count should be zero or greater. Provided:%d", val))
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
	for _, c := range []hal.Capability{hal.DigitalOutput, hal.DigitalInput, hal.PWM, hal.AnalogInput} {
		if v, ok := parameters[cap2string(c)]; ok {
			val, ok := hal.ConvertToInt(v)
			if !ok {
				return nil, fmt.Errorf("failed to type cast '%s' parameter value '%v' as integer", c, v)
			}
			for i := 0; i < val; i++ {
				pins[c] = append(pins[c], i)
			}
		}
	}
	return &driver{
		meta:    f.meta,
		address: parameters[Address].(string),
		pins:    pins,
		client:  f.client,
	}, nil
}
