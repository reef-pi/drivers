package pca9685

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
)

type pcaFactory struct {
	meta       hal.Metadata
	parameters []hal.ConfigParameter
}

var factory *pcaFactory
var once sync.Once

// Factory returns a singleton File based Digital Driver factory
func Factory() hal.DriverFactory {

	once.Do(func() {
		factory = &pcaFactory{
			meta: hal.Metadata{
				Name:        "pca9685",
				Description: "Supports one PCA9685 chip",
				Capabilities: []hal.Capability{
					hal.PWM, hal.DigitalOutput,
				},
			},
			parameters: []hal.ConfigParameter{
				{
					Name:    "Address",
					Type:    hal.Integer,
					Order:   0,
					Default: 0x40,
				},
				{
					Name:    "Frequency",
					Type:    hal.Integer,
					Order:   1,
					Default: 150,
				},
			},
		}
	})

	return factory
}

func (f *pcaFactory) Metadata() hal.Metadata {
	return f.meta
}

func (f *pcaFactory) GetParameters() []hal.ConfigParameter {
	return f.parameters
}

func (f *pcaFactory) ValidateParameters(parameters map[string]interface{}) (bool, map[string][]string) {

	var failures = make(map[string][]string)
	var v interface{}
	var ok bool

	if v, ok = parameters["Address"]; ok {
		val, ok := hal.ConvertToInt(v)
		if !ok {
			failure := fmt.Sprint("Address is not a number. ", v, " was received.")
			failures["Address"] = append(failures["Address"], failure)
		}
		if val <= 0 || val >= 256 {
			failure := fmt.Sprint("Address is out of range (1 - 255). ", v, " was received.")
			failures["Address"] = append(failures["Address"], failure)
		}
	} else {
		failure := fmt.Sprint("Address is required parameter, but was not received.")
		failures["Address"] = append(failures["Address"], failure)
	}

	if v, ok = parameters["Frequency"]; ok {
		val, ok := hal.ConvertToInt(v)
		if !ok {
			failure := fmt.Sprint("Frequency is not a number. ", v, " was received.")
			failures["Frequency"] = append(failures["Frequency"], failure)
		}
		if val <= 0 || val > 1500 {
			failure := fmt.Sprint("Frequency is out of range (1 - 1500). ", v, " was received.")
			failures["Frequency"] = append(failures["Frequency"], failure)
		}
	} else {
		failure := fmt.Sprint("Frequency is required parameter, but was not received.")
		failures["Frequency"] = append(failures["Frequency"], failure)
	}

	return len(failures) == 0, failures
}

func (f *pcaFactory) NewDriver(parameters map[string]interface{}, hardwareResources interface{}) (hal.Driver, error) {
	if valid, failures := f.ValidateParameters(parameters); !valid {
		return nil, errors.New(hal.ToErrorString(failures))
	}

	address, _ := hal.ConvertToInt(parameters["Address"])
	frequency, _ := hal.ConvertToInt(parameters["Frequency"])

	config := PCA9685Config{
		Address:   address,
		Frequency: frequency,
	}

	bus := hardwareResources.(i2c.Bus)

	hwDriver := &PCA9685{
		addr: byte(address),
		bus:  bus,
		Freq: frequency,
	}

	pwm := pca9685Driver{
		mu:       &sync.Mutex{},
		hwDriver: hwDriver,
	}
	if config.Frequency == 0 {
		log.Println("WARNING: pca9685 driver pwm frequency set to 0. Falling back to 1500")
		config.Frequency = 1500
	}
	hwDriver.Freq = config.Frequency // overriding default

	// Create the 16 channels the hardware has
	for i := 0; i < 16; i++ {
		ch := &pca9685Channel{
			channel: i,
			driver:  &pwm,
		}
		pwm.channels = append(pwm.channels, ch)
	}

	// Wake the hardware
	return &pwm, hwDriver.Wake()
}
