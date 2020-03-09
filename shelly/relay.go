package shelly

import (
	"fmt"
	"net/http"
	"time"
)

type HTTPGetter func(string) (*http.Response, error)

type Relay struct {
	channel int
	addr    string
	state   bool
	name    string
	getter  HTTPGetter
}

func NewRelay(name, addr string, channel int, getter HTTPGetter) *Relay {
	r := Relay{
		channel: channel,
		addr:    addr,
		name:    name,
		getter:  getter,
	}
	if getter == nil {
		h := &http.Client{Timeout: 5 * time.Second}
		r.getter = h.Get
	}
	return &r
}

func (r *Relay) Close() error    { return nil }
func (r *Relay) Number() int     { return r.channel }
func (r *Relay) LastState() bool { return r.state }
func (r *Relay) Name() string    { return r.name }

func (r *Relay) Write(b bool) error {
	action := "on"
	if !b {
		action = "off"
	}
	resp, err := r.getter(fmt.Sprintf("%s/relay/%d?turn=%s", r.addr, r.channel, action))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http failure. Code:%d", resp.StatusCode)
	}
	r.state = b
	return nil
}
