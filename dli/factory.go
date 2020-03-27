package dli

import (
	"errors"
	"fmt"
	"github.com/reef-pi/hal"
)

const (
	_user     = "username"
	_password = "password"
	_addr     = "ip"
)

type Factory struct {
	meta       hal.Metadata
	parameters []hal.ConfigParameter
}

func Adapter() hal.DriverFactory {
	return &Factory{
		meta: hal.Metadata{
			Name:         "DLI-Webpowerswitch",
			Description:  "DLI Web Powerswitch pro",
			Capabilities: []hal.Capability{hal.DigitalOutput},
		},
		parameters: []hal.ConfigParameter{
			{
				Name:    _addr,
				Type:    hal.String,
				Order:   0,
				Default: "192.168.1.33",
			},
			{
				Name:    _user,
				Type:    hal.String,
				Order:   1,
				Default: "admin",
			},
			{
				Name:    _password,
				Type:    hal.String,
				Order:   2,
				Default: "1234",
			},
		},
	}
}

func (f *Factory) Metadata() hal.Metadata {
	return f.meta
}
func (f *Factory) GetParameters() []hal.ConfigParameter {
	return f.parameters
}

func (f *Factory) ValidateParameters(parameters map[string]interface{}) (bool, map[string][]string) {

	var failures = make(map[string][]string)

	if v, ok := parameters[_addr]; ok {
		_, ok := v.(string)
		if !ok {
			failure := fmt.Sprint(_addr, " is not a string. ", v, " was received.")
			failures[_addr] = append(failures[_addr], failure)
		}
	} else {
		failure := fmt.Sprint(_addr, " is a required parameter, but was not received.")
		failures[_addr] = append(failures[_addr], failure)
	}

	if v, ok := parameters[_user]; ok {
		_, ok := v.(string)
		if !ok {
			failure := fmt.Sprint(_user, " is not a string. ", v, " was received.")
			failures[_user] = append(failures[_user], failure)
		}
	} else {
		failure := fmt.Sprint(_user, " is a required parameter, but was not received.")
		failures[_user] = append(failures[_user], failure)
	}

	if v, ok := parameters[_password]; ok {
		_, ok := v.(string)
		if !ok {
			failure := fmt.Sprint(_password, " is not a string. ", v, " was received.")
			failures[_password] = append(failures[_password], failure)
		}
	} else {
		failure := fmt.Sprint(_password, " is a required parameter, but was not received.")
		failures[_password] = append(failures[_password], failure)
	}

	return len(failures) == 0, failures
}

func (f *Factory) NewDriver(params map[string]interface{}, _ interface{}) (hal.Driver, error) {
	if valid, failures := f.ValidateParameters(params); !valid {
		return nil, errors.New(hal.ToErrorString(failures))
	}
	addr := params[_addr].(string)
	user := params[_user].(string)
	password := params[_password].(string)
	return NewDriver(addr, user, password), nil
}
