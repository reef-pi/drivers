package ads1x15

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
)

const addressParam = "Address"

var channelAddresses = [4]uint16{configMuxSingle0, configMuxSingle1, configMuxSingle2, configMuxSingle3}
var channelGains = [4]string{"Gain 1", "Gain 2", "Gain 3", "Gain 4"}
var gainOptions = map[string]uint16{
	"2/3": configGainTwoThirds,
	"1":   configGainOne,
	"2":   configGainTwo,
	"4":   configGainFour,
	"8":   configGainEight,
	"16":  configGainSixteen,
}

type ads1X15Factory struct {
	meta       hal.Metadata
	parameters []hal.ConfigParameter
}

func (f *ads1X15Factory) appendParameters() {
	f.parameters = []hal.ConfigParameter{
		{
			Name:    addressParam,
			Type:    hal.Integer,
			Order:   0,
			Default: 0x48,
		},
	}

	for i, name := range channelGains {
		gainParam := hal.ConfigParameter{
			Name:    name,
			Type:    hal.String,
			Order:   i + 1,
			Default: "2/3",
		}
		f.parameters = append(f.parameters, gainParam)
	}
}

func (f *ads1X15Factory) Metadata() hal.Metadata {
	return f.meta
}

func (f *ads1X15Factory) GetParameters() []hal.ConfigParameter {
	return f.parameters
}

func (f *ads1X15Factory) ValidateParameters(parameters map[string]interface{}) (bool, map[string][]string) {

	var failures = make(map[string][]string)
	var v interface{}
	var val interface{}
	var ok bool

	if v, ok = parameters[addressParam]; ok {
		val, ok = hal.ConvertToInt(v)
		if !ok {
			failure := fmt.Sprint(addressParam, " is not a number. ", v, " was received.")
			failures[addressParam] = append(failures[addressParam], failure)
		}
		if val.(int) <= 0 || val.(int) >= 256 {
			failure := fmt.Sprint(addressParam, " is out of range (1 - 255). ", v, " was received.")
			failures[addressParam] = append(failures[addressParam], failure)
		}
	} else {
		failure := fmt.Sprint(addressParam, " is required parameter, but was not received.")
		failures[addressParam] = append(failures[addressParam], failure)
	}

	for _, channelGain := range channelGains {
		if v, ok = parameters[channelGain]; ok {
			if _, err := parseGain(parameters[channelGain]); err != nil {
				failures[channelGain] = append(failures[channelGain], fmt.Sprint(channelGain, err.Error()))
			}
		} else {
			failure := fmt.Sprint(channelGain, " is a required parameter, but was not received.")
			failures[channelGain] = append(failures[channelGain], failure)
		}
	}
	return len(failures) == 0, failures
}

func parseGain(v interface{}) (string, error) {
	var val interface{}
	var ok bool

	val, ok = v.(string)
	if !ok {
		val, ok = hal.ConvertToInt(v)
		if ok {
			val = strconv.Itoa(val.(int))
		}
	}

	if !ok {
		failure := fmt.Sprint(" is not a string. ", v, " was received.")
		return "", fmt.Errorf(failure)
	}

	if _, ok = gainOptions[val.(string)]; !ok {
		failure := fmt.Sprint(" is not a valid value of 2/3, 1, 2, 4, 8, or 16. ", v, " was received.")
		return "", fmt.Errorf(failure)
	}

	return val.(string), nil
}

func (f *ads1X15Factory) newDriver(parameters map[string]interface{}, hardwareResources interface{}, shift int, delay time.Duration) (hal.Driver, error) {
	if valid, failures := f.ValidateParameters(parameters); !valid {
		return nil, errors.New(hal.ToErrorString(failures))
	}

	intAddress, _ := hal.ConvertToInt(parameters[addressParam])
	address := byte(intAddress)
	bus := hardwareResources.(i2c.Bus)

	var configRegister [2]byte
	if err := bus.ReadFromReg(address, 0x01, configRegister[:]); err != nil {
		return nil, err
	}

	var driver = driver{
		meta:     f.meta,
		channels: []hal.AnalogInputPin{},
	}

	// Create the 4 channels the hardware has
	for i, channelAddress := range channelAddresses {
		gain, _ := parseGain(parameters[channelGains[i]])
		ch, err := newChannel(bus, address, i, channelAddress, gainOptions[gain], shift, delay)
		if err != nil {
			return nil, err
		}

		driver.channels = append(driver.channels, ch)
	}

	return &driver, nil
}
