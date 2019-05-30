package pca9685

import (
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
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
	config   PCA9685Config
	hwDriver *PCA9685
	mu       *sync.Mutex
	channels []*pca9685Channel
}

var DefaultPCA9685Config = PCA9685Config{
	Address:   0x40,
	Frequency: 1500,
}

func HALAdpater(config PCA9685Config, bus i2c.Bus) (hal.Driver, error) {

	hwDriver := New(byte(config.Address), bus)
	pwm := pca9685Driver{
		config:   config,
		mu:       &sync.Mutex{},
		hwDriver: hwDriver,
	}
	if config.Frequency == 0 {
		log.Println("WARNING: pca9685 driver pwm frequency set to 0. Falling back to 1500")
		config.Frequency = 1500
	}
	hwDriver.Freq = config.Frequency // overriding default

	// Create the 16 channels the hardware has
	for i := 0; i < 16; i++ {
		ch := &pca9685Channel{
			channel: i,
			driver:  &pwm,
		}
		pwm.channels = append(pwm.channels, ch)
	}

	// Wake the hardware
	return &pwm, hwDriver.Wake()
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
			hal.PWM,
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
func (p *pca9685Driver) OutputPins() []hal.OutputPin {
	pins := make([]hal.OutputPin, len(p.channels))
	for i, ch := range p.channels {
		pins[i] = ch
	}
	return pins
}

func (p *pca9685Driver) OutputPin(n int) (hal.OutputPin, error) {
	return p.PWMChannel(n)
}

// value should be within 0-100
func (p *pca9685Driver) set(pin int, value float64) error {
	if (value > 100) || (value < 0) {
		return fmt.Errorf("invalid pwm range: %f, value should be within 0 to 100", value)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.hwDriver.SetPwm(pin, 0, uint16(value*4095))
}
