package tasmota

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/reef-pi/hal"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// parseOutputs parses an output configuration string and returns a sorted slice of output numbers.
// Supported formats:
// - "1" -> [1]
// - "1,2,3" -> [1, 2, 3]
// - "1-3" -> [1, 2, 3]
// - "1-3,5,7-9" -> [1, 2, 3, 5, 7, 8, 9]
// Returns error for invalid formats, negative numbers, duplicates, or reversed ranges.
func parseOutputs(config string) ([]int, error) {
	if config == "" {
		return nil, errors.New("output configuration cannot be empty")
	}

	outputMap := make(map[int]bool)
	var outputs []int

	// Split by comma
	parts := strings.Split(config, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)

		if strings.Contains(part, "-") {
			// Handle range format (e.g., "1-3")
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range format: %s", part)
			}

			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("invalid start number in range '%s': %v", part, err)
			}

			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid end number in range '%s': %v", part, err)
			}

			if start < 0 || end < 0 {
				return nil, fmt.Errorf("output numbers must be non-negative, got range '%s'", part)
			}

			if start > end {
				return nil, fmt.Errorf("invalid range '%s': start (%d) is greater than end (%d)", part, start, end)
			}

			for i := start; i <= end; i++ {
				if outputMap[i] {
					return nil, fmt.Errorf("duplicate output number: %d", i)
				}
				outputMap[i] = true
				outputs = append(outputs, i)
			}
		} else {
			// Handle single number format
			num, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid output number '%s': %v", part, err)
			}

			if num < 0 {
				return nil, fmt.Errorf("output numbers must be non-negative, got %d", num)
			}

			if outputMap[num] {
				return nil, fmt.Errorf("duplicate output number: %d", num)
			}

			outputMap[num] = true
			outputs = append(outputs, num)
		}
	}

	if len(outputs) == 0 {
		return nil, errors.New("no valid output numbers found")
	}

	// Sort outputs for consistent ordering
	sort.Ints(outputs)

	return outputs, nil
}

type httpDriver struct {
	meta     hal.Metadata
	address  string
	outputs  []int
	pins     []hal.DigitalOutputPin
	channels []hal.PWMChannel
}

// pinDriver represents a digital output pin on a Tasmota device
type pinDriver struct {
	driver *httpDriver
	number int
}

// channelDriver represents a PWM channel on a Tasmota device
type channelDriver struct {
	driver *httpDriver
	number int
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
	return m.channels
}

func (m *httpDriver) PWMChannel(index int) (hal.PWMChannel, error) {
	if index < 0 || index >= len(m.channels) {
		return nil, fmt.Errorf("PWM channel index %d out of range (0-%d)", index, len(m.channels)-1)
	}
	return m.channels[index], nil
}

func (m *httpDriver) doRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	c := http.Client{
		Timeout: 5 * time.Second,
	}
	return c.Do(req)
}

func (m *httpDriver) readBody(body io.ReadCloser) ([]byte, error) {
	defer body.Close()
	msg, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return msg, nil
}


func (m *httpDriver) DigitalOutputPins() []hal.DigitalOutputPin {
	return m.pins
}

func (m *httpDriver) DigitalOutputPin(index int) (hal.DigitalOutputPin, error) {
	if index < 0 || index >= len(m.pins) {
		return nil, fmt.Errorf("digital output pin index %d out of range (0-%d)", index, len(m.pins)-1)
	}
	return m.pins[index], nil
}

// pinDriver methods

func (p *pinDriver) Close() error {
	return nil
}

func (p *pinDriver) Name() string {
	return "Tasmota"
}

func (p *pinDriver) Number() int {
	return 0
}

func (p *pinDriver) Write(b bool) error {
	const baseUri = "http://%s/cm?cmnd=Power%d%%20%t"
	uri := fmt.Sprintf(baseUri, p.driver.address, p.number, b)
	resp, err := p.driver.doRequest(uri)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		return nil
	}
	body, err := p.driver.readBody(resp.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf("HTTP Code:%d. Body:%v", resp.StatusCode, string(body))
}

func (p *pinDriver) LastState() bool {
	const urlBase = "http://%s/cm?cmnd=Power%d"
	uri := fmt.Sprintf(urlBase, p.driver.address, p.number)
	resp, err := p.driver.doRequest(uri)
	if err != nil {
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}
	body, err := p.driver.readBody(resp.Body)
	if err != nil {
		return false
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return false
	}

	if result[fmt.Sprintf("POWER%d", p.number)] == "ON" {
		return true
	}

	if result["POWER"] == "ON" {
		return true
	}

	return false
}

// channelDriver methods

func (c *channelDriver) Close() error {
	return nil
}

func (c *channelDriver) Name() string {
	return "Tasmota"
}

func (c *channelDriver) Number() int {
	return 0
}

func (c *channelDriver) Set(value float64) error {
	const urlBase = "http://%s/cm?cmnd=Dimmer%%20%.0f"
	uri := fmt.Sprintf(urlBase, c.driver.address, value)
	resp, err := c.driver.doRequest(uri)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		return nil
	}
	body, err := c.driver.readBody(resp.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf("HTTP Code:%d. Body:%v", resp.StatusCode, string(body))
}

func (c *channelDriver) Write(b bool) error {
	const baseUri = "http://%s/cm?cmnd=Power%d%%20%t"
	uri := fmt.Sprintf(baseUri, c.driver.address, c.number, b)
	resp, err := c.driver.doRequest(uri)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		return nil
	}
	body, err := c.driver.readBody(resp.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf("HTTP Code:%d. Body:%v", resp.StatusCode, string(body))
}

func (c *channelDriver) LastState() bool {
	const urlBase = "http://%s/cm?cmnd=Power%d"
	uri := fmt.Sprintf(urlBase, c.driver.address, c.number)
	resp, err := c.driver.doRequest(uri)
	if err != nil {
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}
	body, err := c.driver.readBody(resp.Body)
	if err != nil {
		return false
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return false
	}

	if result[fmt.Sprintf("POWER%d", c.number)] == "ON" {
		return true
	}

	if result["POWER"] == "ON" {
		return true
	}

	return false
}

type factory struct {
	meta       hal.Metadata
	parameters []hal.ConfigParameter
}

var pwmDriverFactory *factory
var once sync.Once

const address = "Address"
const output = "Output"

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
				{
					Name:    output,
					Type:    hal.Integer,
					Order:   1,
					Default: 0,
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

	if v, ok := parameters[output]; ok {
		// Accept both string and integer for backward compatibility
		var outputConfig string
		switch val := v.(type) {
		case string:
			outputConfig = val
		case int:
			outputConfig = strconv.Itoa(val)
		default:
			failure := fmt.Sprint(output, " must be string or integer. ", v, " was received.")
			failures[output] = append(failures[output], failure)
		}

		if len(outputConfig) > 0 {
			// Validate the output configuration
			_, err := parseOutputs(outputConfig)
			if err != nil {
				failure := fmt.Sprint(output, " configuration is invalid: ", err.Error())
				failures[output] = append(failures[output], failure)
			}
		}
	} else {
		failure := fmt.Sprint(output, " is a required parameter, but was not received.")
		failures[output] = append(failures[output], failure)
	}

	return len(failures) == 0, failures
}

func (f *factory) Metadata() hal.Metadata {
	return f.meta
}

func (f *factory) NewDriver(parameters map[string]interface{}, hardwareResources interface{}) (hal.Driver, error) {
	if parameters[output] == nil {
		parameters[output] = "1"
	}

	// Convert output to string for consistent processing
	var outputStr string
	if outputInt, ok := parameters[output].(int); ok {
		outputStr = strconv.Itoa(outputInt)
	} else if str, ok := parameters[output].(string); ok {
		outputStr = str
	} else {
		return nil, fmt.Errorf("output must be int or string, got %T", parameters[output])
	}
	parameters[output] = outputStr

	if valid, failures := f.ValidateParameters(parameters); !valid {
		return nil, errors.New(hal.ToErrorString(failures))
	}

	// Parse output configuration
	outputs, err := parseOutputs(outputStr)
	if err != nil {
		return nil, fmt.Errorf("invalid output configuration: %v", err)
	}

	// Create driver
	driver := &httpDriver{
		meta:     f.meta,
		address:  parameters[address].(string),
		outputs:  outputs,
		pins:     []hal.DigitalOutputPin{},
		channels: []hal.PWMChannel{},
	}

	// Create pin and channel objects for each output
	for _, outNum := range outputs {
		driver.pins = append(driver.pins, &pinDriver{
			driver: driver,
			number: outNum,
		})
		driver.channels = append(driver.channels, &channelDriver{
			driver: driver,
			number: outNum,
		})
	}

	return driver, nil
}
