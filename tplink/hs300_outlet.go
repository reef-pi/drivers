package tplink

import (
	"encoding/json"
	"fmt"

	"github.com/reef-pi/hal"
)

type (
	Outlet struct {
		name       string
		id         string
		addr       string
		state      bool
		cnFactory  ConnectionFactory
		calibrator hal.Calibrator
	}
)

func (o *Outlet) Name() string {
	return o.name
}

func (o *Outlet) Write(state bool) error {
	if state {
		return o.On()
	}
	return o.Off()
}

func (o *Outlet) RTEmeter() (*HS300Realtime, error) {
	var cmd HS300EmeterCmd
	cmd.Context.Children = []string{o.id}
	d, err := command(o.cnFactory, o.addr, &cmd)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(d, &cmd); err != nil {
		return nil, err
	}
	return &cmd.Emeter.Realtime, nil
}

func (o *Outlet) LastState() bool {
	return o.state
}

func (o *Outlet) On() error {
	cmd := new(CmdRelayState)
	cmd.System.RelayState.State = 1
	cmd.Context.Children = []string{o.id}
	if _, err := command(o.cnFactory, o.addr, cmd); err != nil {
		return err
	}
	o.state = true
	return nil
}
func (o *Outlet) Off() error {
	cmd := new(CmdRelayState)
	cmd.System.RelayState.State = 0
	cmd.Context.Children = []string{o.id}
	if _, err := command(o.cnFactory, o.addr, cmd); err != nil {
		return err
	}
	o.state = true
	return nil
}
func (o *Outlet) Read() (float64, error) {
	em, err := o.RTEmeter()
	if err != nil {
		return 0, err
	}
	return em.Current, nil
}

func (o *Outlet) Calibrate(points []hal.Measurement) error {
	cal, err := hal.CalibratorFactory(points)
	if err != nil {
		return err
	}
	o.calibrator = cal
	return nil
}
func (o *Outlet) Measure() (float64, error) {
	v, err := o.Read()
	if err != nil {
		return 0, err
	}
	if o.calibrator == nil {
		return 0, fmt.Errorf("Not calibrated")
	}
	return o.calibrator.Calibrate(v), nil
}

func (o *Outlet) Close() error {
	return nil
}
