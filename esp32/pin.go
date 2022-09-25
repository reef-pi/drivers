package esp32

import (
	"errors"
	"fmt"
	"github.com/reef-pi/hal"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var _notImplemented = errors.New("not implemented")

const _true = "true"

var ErrIncompatibleCapability = errors.New("incompatible capability")

type pin struct {
	address string
	number  int
	cap     hal.Capability
	client  HTTPClient
}

func (p *pin) Close() error {
	return nil
}

func (p *pin) Number() int {
	return p.number
}
func (p *pin) Name() string {
	return fmt.Sprintf("capability:%s pin:%d", p.cap.String(), p.number)
}

func (p *pin) Calibrate([]hal.Measurement) error {
	return _notImplemented
}

func (p *pin) Measure() (float64, error) {
	return 0, _notImplemented
}

func (p *pin) doRequest(verb, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(verb, url, body)
	if err != nil {
		return nil, err
	}
	return p.client(req)
}

func (p *pin) readBody(body io.ReadCloser) ([]byte, error) {
	defer body.Close()
	msg, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (p *pin) LastState() bool {
	baseUri := "http://%s/outlets/%d"
	uri := fmt.Sprintf(baseUri, p.address, p.number)
	resp, err := p.doRequest(http.MethodGet, uri, nil)
	if err != nil {
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}
	_, err = p.readBody(resp.Body)
	if err != nil {
		return false
	}
	return true
}

func mapTo255String(f float64) string {
	if f < 0 {
		f = 0
	}
	if f > 100 {
		f = 100
	}
	return strconv.Itoa(int(f * 255 / 100))
}

func (p *pin) incompatibleCapability() error {
	return fmt.Errorf("%w. supported capability:%s", ErrIncompatibleCapability, p.cap.String())
}

func (p *pin) Set(v float64) error {
	if p.cap != hal.PWM {
		return p.incompatibleCapability()
	}
	baseUri := "http://%s/jacks/%d"
	uri := fmt.Sprintf(baseUri, p.address, p.number)
	resp, err := p.doRequest(http.MethodPost, uri, strings.NewReader(mapTo255String(v)))

	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		return nil
	}
	body, err := p.readBody(resp.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf("HTTP Code:%d. Body:%v", resp.StatusCode, string(body))
}

func (p *pin) Write(b bool) error {
	if p.cap != hal.DigitalOutput {
		return p.incompatibleCapability()
	}
	action := "off"
	if b {
		action = "off"
	}
	baseUri := "http://%s/outlets/%d/%s"
	uri := fmt.Sprintf(baseUri, p.address, p.number, action)
	resp, err := p.doRequest(http.MethodPost, uri, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		return nil
	}
	body, err := p.readBody(resp.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf("HTTP Code:%d. Body:%v", resp.StatusCode, string(body))
}

func (p *pin) Read() (bool, error) {
	if p.cap != hal.DigitalInput {
		return false, p.incompatibleCapability()
	}
	baseUri := "http://%s/inlets/%d"
	uri := fmt.Sprintf(baseUri, p.address, p.number)
	resp, err := p.doRequest(http.MethodGet, uri, nil)
	if err != nil {
		return false, err
	}
	body, err := p.readBody(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read http response body. Error:%w", err)
	}
	if resp.StatusCode != 200 {
		return false, fmt.Errorf("HTTP Code:%d. Body:%v", resp.StatusCode, string(body))
	}
	return strings.ToLower(strings.TrimSpace(string(body))) == _true, nil
}

func (p *pin) Value() (float64, error) {
	if p.cap != hal.AnalogInput {
		return 0, p.incompatibleCapability()
	}
	baseUri := "http://%s/analog_inputs/%d"
	uri := fmt.Sprintf(baseUri, p.address, p.number)
	resp, err := p.doRequest(http.MethodGet, uri, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to make http request. Error:%w", err)
	}
	body, err := p.readBody(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read http response body. Error:%w", err)
	}
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("HTTP Code:%d. Body:%v", resp.StatusCode, string(body))
	}
	return strconv.ParseFloat(strings.TrimSpace(string(body)), 64)
}
