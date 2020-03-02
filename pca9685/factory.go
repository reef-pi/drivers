package pca9685

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
)

const addressParam = "Address"
const freqParam = "Frequency"

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
					Name:    addressParam,
					Type:    hal.Integer,
					Order:   0,
					Default: 0x40,
				},
				{
					Name:    freqParam,
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

	if v, ok = parameters[addressParam]; ok {
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
		failure := fmt.Sprint(addressParam, " is required parameter, but was not received.")
		failures[addressParam] = append(failures[addressParam], failure)
	}

	if v, ok = parameters[freqParam]; ok {
		val, ok := hal.ConvertToInt(v)
		if !ok {
			failure := fmt.Sprint(freqParam, " is not a number. ", v, " was received.")
			failures[freqParam] = append(failures[freqParam], failure)
		}
		if val <= 0 || val > 1500 {
			failure := fmt.Sprint(freqParam, " is out of range (1 - 1500). ", v, " was received.")
			failures[freqParam] = append(failures[freqParam], failure)
		}
	} else {
		failure := fmt.Sprint(freqParam, " is required parameter, but was not received.")
		failures[freqParam] = append(failures[freqParam], failure)
	}

	return len(failures) == 0, failures
}

func (f *pcaFactory) NewDriver(parameters map[string]interface{}, hardwareResources interface{}) (hal.Driver, error) {
	if valid, failures := f.ValidateParameters(parameters); !valid {
		return nil, errors.New(hal.ToErrorString(failures))
	}

	address, _ := hal.ConvertToInt(parameters[addressParam])
	frequency, _ := hal.ConvertToInt(parameters[freqParam])

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
