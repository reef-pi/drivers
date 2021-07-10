package tasmota

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/reef-pi/hal"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type httpDriver struct {
	meta    hal.Metadata
	address string
}

func (m *httpDriver) Close() error {
	return nil
}

func (m *httpDriver) Metadata() hal.Metadata {
	return m.meta
}

func (m *httpDriver) Name() string {
	return "Tasmota"
}

func (m *httpDriver) Number() int {
	return 0
}

func (m *httpDriver) Pins(capability hal.Capability) ([]hal.Pin, error) {
	switch capability {
	case hal.DigitalOutput:
		return []hal.Pin{m}, nil
	case hal.PWM:
		return []hal.Pin{m}, nil
	default:
		return nil, fmt.Errorf("unsupported capability:%s", capability.String())
	}
}

func (m *httpDriver) PWMChannels() []hal.PWMChannel {
	return []hal.PWMChannel{m}
}

func (m *httpDriver) PWMChannel(_ int) (hal.PWMChannel, error) {
	return m, nil
}

func (m *httpDriver) LastState() bool {
	uri := fmt.Sprintf("http://%s/cm?cmnd=Power0", m.address)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return false
	}
	c := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := c.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	msg, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return false
	}
	var result map[string]interface{}
	err = json.Unmarshal([]byte(msg), &result)
	if err != nil {
		return false
	}
	return result["POWER"] == "ON"
}

func (m *httpDriver) Set(value float64) error {
	uri := fmt.Sprintf("http://%s/cm?cmnd=Dimmer%%20%.0f", m.address, value)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}
	c := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	msg, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		return nil
	}
	return fmt.Errorf("HTTP Code:%d. Body:%v", resp.StatusCode, string(msg))
}

func (m *httpDriver) Write(b bool) error {
	uri := fmt.Sprintf("http://%s/cm?cmnd=Power0%%20%t", m.address, b)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}
	c := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	msg, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		return nil
	}
	return fmt.Errorf("HTTP Code:%d. Body:%v", resp.StatusCode, string(msg))
}

func (m *httpDriver) DigitalOutputPins() []hal.DigitalOutputPin {
	return []hal.DigitalOutputPin{m}
}

func (m *httpDriver) DigitalOutputPin(_ int) (hal.DigitalOutputPin, error) {
	return m, nil
}

type factory struct {
	meta       hal.Metadata
	parameters []hal.ConfigParameter
}

var pwmDriverFactory *factory
var once sync.Once

const address = "Domain or Address"

func HttpDriverFactory() hal.DriverFactory {

	once.Do(func() {
		pwmDriverFactory = &factory{
			meta: hal.Metadata{
				Name:         "Tasmota Http",
				Description:  "Tasmota Http Driver",
				Capabilities: []hal.Capability{hal.PWM, hal.DigitalOutput},
			},
			parameters: []hal.ConfigParameter{
				{
					Name:    address,
					Type:    hal.String,
					Order:   0,
					Default: "192.1.168.4",
				},
			},
		}
	})

	return pwmDriverFactory
}

func (f *factory) GetParameters() []hal.ConfigParameter {
	return f.parameters
}

func (f *factory) ValidateParameters(parameters map[string]interface{}) (bool, map[string][]string) {
	var failures = make(map[string][]string)

	if v, ok := parameters[address]; ok {
		val, ok := v.(string)
		if !ok {
			failure := fmt.Sprint(address, " is not a string. ", v, " was received.")
			failures[address] = append(failures[address], failure)
		} else if len(val) <= 0 {
			failure := fmt.Sprint(address, " empty values are not allowed.")
			failures[address] = append(failures[address], failure)
		} else if len(val) >= 256 {
			failure := fmt.Sprint(address, " size should be lower than 255 characters. ", val, " was received.")
			failures[address] = append(failures[address], failure)
		}
	} else {
		failure := fmt.Sprint(address, " is a required parameter, but was not received.")
		failures[address] = append(failures[address], failure)
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
	driver := &httpDriver{
		meta: f.meta,
		address: parameters[address].(string),
	}
	return driver, nil
}