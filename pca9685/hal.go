package pca9685

import (
	"fmt"
	"sort"
	"sync"

	"github.com/reef-pi/hal"
)

type PCA9685Config struct {
	Address   int `json:"address"` // 0x40
	Frequency int `json:"frequency"`
}

type pca9685Channel struct {
	driver  *pca9685Driver
	channel int
	v       float64
}

func (c *pca9685Channel) Name() string { return fmt.Sprintf("%d", c.channel) }
func (c *pca9685Channel) Number() int  { return c.channel }
func (c *pca9685Channel) Close() error { return nil }
func (c *pca9685Channel) Set(value float64) error {
	if err := c.driver.set(c.channel, value); err != nil {
		return err
	}
	c.v = value
	return nil
}
func (c *pca9685Channel) Write(b bool) error {
	var v float64
	if b {
		v = 100
	}
	if err := c.driver.set(c.channel, v); err != nil {
		return err
	}
	c.v = v
	return nil
}

func (c *pca9685Channel) LastState() bool { return c.v == 100 }

type pca9685Driver struct {
	hwDriver *PCA9685
	mu       *sync.Mutex
	channels []*pca9685Channel
}

func (p *pca9685Driver) Close() error {
	// Close the driver (will clear all registers)
	if err := p.hwDriver.Close(); err != nil {
		return err
	}
	// Send the hardware to sleep
	return p.hwDriver.Sleep()
}

func (p *pca9685Driver) Metadata() hal.Metadata {
	return hal.Metadata{
		Name:        "pca9685",
		Description: "Supports one PCA9685 chip",
		Capabilities: []hal.Capability{
			hal.PWM, hal.DigitalOutput,
		},
	}
}

func (p *pca9685Driver) PWMChannels() []hal.PWMChannel {
	// Return array of channels soreted by name
	var chs []hal.PWMChannel
	for _, ch := range p.channels {
		chs = append(chs, ch)
	}
	sort.Slice(chs, func(i, j int) bool { return chs[i].Name() < chs[j].Name() })
	return chs
}

func (p *pca9685Driver) PWMChannel(chnum int) (hal.PWMChannel, error) {
	// Return given channel
	if chnum < 0 || chnum >= len(p.channels) {
		return nil, fmt.Errorf("invalid channel %d", chnum)
	}
	return p.channels[chnum], nil
}
func (p *pca9685Driver) DigitalOutputPins() []hal.DigitalOutputPin {
	pins := make([]hal.DigitalOutputPin, len(p.channels))
	for i, ch := range p.channels {
		pins[i] = ch
	}
	return pins
}

func (p *pca9685Driver) DigitalOutputPin(n int) (hal.DigitalOutputPin, error) {
	return p.PWMChannel(n)
}

// value should be within 0-100
func (p *pca9685Driver) set(pin int, value float64) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	switch {
	case value > 100:
		return fmt.Errorf("invalid value: %f above 100", value)
	case value < 0:
		return fmt.Errorf("invalid value: %f below 0", value)
	case value == 0:
		return p.hwDriver.SetPwm(pin, 0, 4096)
	case value == 100:
		return p.hwDriver.SetPwm(pin, 4096, 0)
	default:
		return p.hwDriver.SetPwm(pin, 0, uint16(value*40.95))
	}
}

func (p *pca9685Driver) Pins(cap hal.Capability) ([]hal.Pin, error) {
	switch cap {
	case hal.DigitalOutput, hal.PWM:
		var pins []hal.Pin
		for _, pin := range p.channels {
			pins = append(pins, pin)
		}
		return pins, nil
	default:
		return nil, fmt.Errorf("unsupported capability: %s", cap.String())
	}
}
