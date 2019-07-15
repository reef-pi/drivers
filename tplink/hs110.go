package tplink

import (
	"encoding/json"
	"github.com/reef-pi/hal"
)

type (
	EmeterCmd struct {
		Emeter struct {
			Realtime Realtime `json:"get_realtime"`
		} `json:"emeter"`
	}
	Realtime struct {
		Current  float64 `json:"current,omitempty"`
		Voltage  float64 `json:"voltage,omitempty"`
		Power    float64 `json:"power,omitempty"`
		Total    float64 `json:"total,omitempty"`
		ErrrCode int     `json:"err_code,omitempty"`
	}
)

type HS110Plug struct {
	HS103Plug
}

func NewHS110Plug(addr string) *HS110Plug {
	return &HS110Plug{
		HS103Plug{
			addr: addr,
			meta: hal.Metadata{
				Name:        "tplink-hs110",
				Description: "tplink hs110 series smart plug driver with current monitoring",
				Capabilities: []hal.Capability{
					hal.Output,
				},
			},
			cnFactory: TCPConnFactory,
		},
	}
}

func (p *HS110Plug) RTEmeter() (*Realtime, error) {
	d, err := command(p.cnFactory, p.addr, new(EmeterCmd))
	if err != nil {
		return nil, err
	}
	var cmd EmeterCmd
	if err := json.Unmarshal(d, &cmd); err != nil {
		return nil, err
	}
	return &cmd.Emeter.Realtime, nil
}
