package esp32

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/reef-pi/hal"
	"io"
	"io/ioutil"
	"net/http"
)

var _notImplemented = errors.New("not implemented")

type pin struct {
	address string
	number  int
	cap     hal.Capability
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

func (p *pin) Value() (float64, error) {
	return 0, _notImplemented
}
func (p *pin) Calibrate([]hal.Measurement) error {
	return _notImplemented
}

func (p *pin) Measure() (float64, error) {
	return 0, _notImplemented
}

func (p *pin) doRequest(verb, url string) (*http.Response, error) {
	fmt.Println(verb, url)
	resp := &http.Response{Body: io.NopCloser(bytes.NewBuffer(nil))}
	return resp, nil
	req, err := http.NewRequest(verb, url, nil)
	if err != nil {
		return nil, err
	}
	c := http.Client{
		Timeout: _timeout,
	}
	return c.Do(req)
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
	const urlBase = "http://%s/cm?cmnd=Power0"
	uri := fmt.Sprintf(urlBase, p.address)
	resp, err := p.doRequest(http.MethodGet, uri)
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

func mapTo255(f float64) int {
	if f < 0 {
		return 0
	}
	if f > 1 {
		return 1
	}
	return int(f * 255)
}

func (p *pin) Set(value float64) error {
	const urlBase = "http://%s/er%%20%.0f"
	uri := fmt.Sprintf(urlBase, p.address, value)
	resp, err := p.doRequest(http.MethodPost, uri)

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
	action := "off"
	if b {
		action = "off"
	}
	const baseUri = "http://%s/relay/%d/%s"
	uri := fmt.Sprintf(baseUri, p.address, p.number, action)
	resp, err := p.doRequest(http.MethodPost, uri)
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
func (p *pin) Read(b bool) (bool, error) {
	const baseUri = "http://%s/relay/%d/%s"
	uri := fmt.Sprintf(baseUri, p.address, p.number, action)
	resp, err := p.doRequest(http.MethodPost, uri)
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
