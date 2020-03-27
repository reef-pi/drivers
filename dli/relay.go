package dli

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Relay struct {
	channel int
	config  Config
	state   bool
}

func (r *Relay) Close() error    { return nil }
func (r *Relay) Number() int     { return r.channel }
func (r *Relay) LastState() bool { return r.state }
func (r *Relay) Name() string    { return fmt.Sprintf("DLI-webpowerswitch-pro-%d", r.channel) }

func (r *Relay) Write(state bool) error {
	uri := fmt.Sprintf("http://%s/restapi/relay/outlets/%d/state/", r.config.addr, r.channel)
	req, err := http.NewRequest("PUT", uri, bytes.NewBuffer([]byte("value=true")))
	if err != nil {
		return err
	}
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	c := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	v := url.Values{}
	if state {
		v.Add("value", "true")
	} else {
		v.Add("value", "false")
	}
	req, err = http.NewRequest("PUT", uri, strings.NewReader(v.Encode()))
	r.config.setDigestAuth(req, resp)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err = c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 204 || resp.StatusCode == 200 {
		r.state = state
		return nil
	}
	msg, _ := ioutil.ReadAll(resp.Body)
	return fmt.Errorf("HTTP Code:%d. Body:%v", resp.StatusCode, string(msg))
}
