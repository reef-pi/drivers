package pca9685

import (
	"testing"

	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
)

var conf = []byte(`{"address":66, "frequency":200}`)

func TestHALAdapter(t *testing.T) {
	driver, err := HALAdapter(conf, i2c.MockBus())
	if err != nil {
		t.Errorf("unexpected error making driver %v", err)
	}

	meta := driver.Metadata()
	if meta.Name != "pca9685" {
		t.Errorf("name %s did not match pca9685", meta.Name)
	}
	if meta.Capabilities[0] != hal.PWM {
		t.Error("driver didn't indicate it supports PWM")
	}
	if meta.Capabilities[1] != hal.DigitalOutput {
		t.Error("driver didn't indicate it supports digital output")
	}
	pwmDriver, ok := driver.(hal.PWMDriver)
	if !ok {
		t.Error("driver is not a PWM interface")
	}
	if pwmDriver == nil {
		t.Error("driver is nil")
	}

	channels := pwmDriver.PWMChannels()
	if l := len(channels); l != 16 {
		t.Errorf("expected 16 channels, got %d", l)
	}

	channel15, err := pwmDriver.PWMChannel(15)
	if err != nil {
		t.Errorf("error fetching channel 15 %v", err)
	}
	if channel15 == nil {
		t.Error("nil pwm driver")
	}
}

func TestPca9685Channel_Set(t *testing.T) {
	driver, err := HALAdapter(conf, i2c.MockBus())
	if err != nil {
		t.Errorf("unexpected error making driver %v", err)
	}

	pwmDriver := driver.(hal.PWMDriver)
	channel15, err := pwmDriver.PWMChannel(15)
	if err != nil {
		t.Errorf("error fetching channel 15 %v", err)
	}

	err = channel15.Set(50)
	if err != nil {
		t.Error("can't set channel to 50%")
	}

	err = channel15.Set(150)
	if err == nil {
		t.Error("channel 15 allowed setting 150%")
	}

	err = channel15.Set(-1)
	if err == nil {
		t.Error("channel 15 allowed setting -1%")
	}
}

func TestPca9685Driver_Close(t *testing.T) {
	driver, err := HALAdapter(conf, i2c.MockBus())
	if err != nil {
		t.Errorf("unexpected error making driver %v", err)
	}

	err = driver.Close()
	if err != nil {
		t.Errorf("unexpected error closing driver %v", err)
	}
}
