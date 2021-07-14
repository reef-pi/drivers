package tasmota

import (
	"github.com/reef-pi/hal"
	"os"
	"testing"
)

func TestHttpDriver_AsDigitalOut(t *testing.T) {

	address := os.Getenv("TASMOTA_TEST_ADDRESS")

	if 	len(address) == 0 {
		address = "192.168.1.46"
	}

	f := HttpDriverFactory()

	params := map[string]interface{}{
		"Domain or Address": address,
	}

	d, err := f.NewDriver(params, nil)
	if err != nil {
		t.Fatal(err)
	}

	meta := d.Metadata()
	if len(meta.Capabilities) != 2 {
		t.Error("Expected 1 capabilities, found:", len(meta.Capabilities))
	}

	o, ok := d.(hal.DigitalOutputDriver)
	if !ok {
		t.Error("Failed to type driver to Digital output driver")
	}

	if len(o.DigitalOutputPins()) != 1 {
		t.Error("Expected a single digital output pwm pin, found:", len(o.DigitalOutputPins()))
	}

	p, err := o.DigitalOutputPin(0)
	if err != nil || p == nil {
		t.Error("Expected a digital output pin")
	}

	if p.Name() != "Tasmota" {
		t.Error("Expected Tasmota name, found: ", p.Name())
	}

	if p.Number() != 0 {
		t.Error("Expected number 0, found: ", p.Number())
	}

	testRealDevice := os.Getenv("TASMOTA_TEST_REAL_DEVICE")

	if testRealDevice == "True" {

		err = p.Write(true)
		if err != nil {
			t.Error("Expected write true inn the digital output, error: ", err.Error())
		}

		if !p.LastState() {
			t.Error("Expected last state is true")
		}

		err = p.Write(false)
		if err != nil {
			t.Error("Expected write false inn the digital output, error: ", err.Error())
		}

		if p.LastState() {
			t.Error("Expected last state is false")
		}
	}

}

func TestHttpDriver_AsPWMDriver(t *testing.T) {

	address := os.Getenv("TASMOTA_TEST_ADDRESS")

	if 	len(address) == 0 {
		address = "192.168.1.46"
	}

	f := HttpDriverFactory()

	params := map[string]interface{}{
		"Domain or Address": address,
	}

	d, err := f.NewDriver(params, nil)
	if err != nil {
		t.Fatal(err)
	}

	meta := d.Metadata()
	if len(meta.Capabilities) != 2 {
		t.Error("Expected 1 capabilities, found:", len(meta.Capabilities))
	}

	pwm, ok := d.(hal.PWMDriver)
	if !ok {
		t.Error("Failed to type driver to PWM driver")
	}

	if len(pwm.PWMChannels()) != 1 {
		t.Error("Expected a single pwm channel, found:", len(pwm.PWMChannels()))
	}

	p, err := pwm.PWMChannel(0)
	if err != nil {
		t.Error("Expected a pwm pin")
	}

	if p.Name() != "Tasmota" {
		t.Error("Expected Tasmota name, found: ", p.Name())
	}

	if p.Number() != 0 {
		t.Error("Expected number 0, found: ", p.Number())
	}

	testRealDevice := os.Getenv("TASMOTA_TEST_REAL_DEVICE")

	if testRealDevice == "True" {

		err = p.Set(100)
		if err != nil {
			t.Error("Expected to set 100 in the pwm output, error: ", err.Error())
		}

		if !p.LastState() {
			t.Error("Expected last state is true")
		}

		err = p.Set(0)
		if err != nil {
			t.Error("Expected to set 0 in the pwm output, error: ", err.Error())
		}

		if p.LastState() {
			t.Error("Expected last state is false")
		}
	}

}

func TestHttpDriver_FactoryValidateParameters(t *testing.T) {

	f := HttpDriverFactory()

	params := map[string]interface{}{
		"Domain or Address": "192.168.1.46",
	}

	_, err := f.NewDriver(params, nil)
	if err != nil {
		t.Fatal(err)
	}

	params = map[string]interface{}{
		"Domain or Address": "",
	}

	_, err = f.NewDriver(params, nil)
	if err == nil {
		t.Fatal("Expected error")
	}

	params = map[string]interface{}{
		"Domain or Address": 1,
	}

	_, err = f.NewDriver(params, nil)
	if err == nil {
		t.Fatal("Expected error")
	}

	params = map[string]interface{}{
		"Domain or Address": nil,
	}

	_, err = f.NewDriver(params, nil)
	if err == nil {
		t.Fatal("Expected error")
	}

}
