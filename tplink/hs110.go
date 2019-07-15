package tplink

import (
	"encoding/json"
	"fmt"
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
	calibrator hal.Calibrator
}

func NewHS110Plug(addr string) *HS110Plug {
	cal, _ := hal.CalibratorFactory([]hal.Measurement{})
	return &HS110Plug{
		HS103Plug: HS103Plug{
			addr: addr,
			meta: hal.Metadata{
				Name:        "tplink-hs110",
				Description: "tplink hs110 series smart plug driver with current monitoring",
				Capabilities: []hal.Capability{
					hal.Output, hal.PH,
				},
			},
			cnFactory: TCPConnFactory,
		},
		calibrator: cal,
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

func (p *HS110Plug) ADCChannels() []hal.ADCChannel {
	return []hal.ADCChannel{p}
}

func (p *HS110Plug) ADCChannel(i int) (hal.ADCChannel, error) {
	if i != 0 {
		return nil, fmt.Errorf("invalid channel number: %d", i)
	}
	return p, nil
}

func (p *HS110Plug) Read() (float64, error) {
	em, err := p.RTEmeter()
	if err != nil {
		return 0, err
	}
	return em.Current, nil
}

func (p *HS110Plug) Calibrate(points []hal.Measurement) error {
	cal, err := hal.CalibratorFactory(points)
	if err != nil {
		return err
	}
	p.calibrator = cal
	return nil
}
func (p *HS110Plug) Measure() (float64, error) {
	v, err := p.Read()
	if err != nil {
		return 0, err
	}
	if p.calibrator == nil {
		return 0, fmt.Errorf("Not calibrated")
	}
	return p.calibrator.Calibrate(v), nil
}
