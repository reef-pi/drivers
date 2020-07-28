package ads1x15

import (
	"testing"

	"github.com/reef-pi/hal"
)

var params = map[string]interface{}{
	"Address": 72,
	"Gain 1":  "2/3",
	"Gain 2":  "1",
	"Gain 3":  "2",
	"Gain 4":  4,
}

type mock struct {
	Bytes []byte
}

func (m *mock) SetAddress(_ byte) error                      { return nil }
func (m *mock) ReadBytes(addr byte, num int) ([]byte, error) { return m.Bytes, nil }
func (m *mock) WriteBytes(addr byte, value []byte) error     { return nil }
func (m *mock) ReadFromReg(addr, reg byte, value []byte) error {
	if len(m.Bytes) >= 2 {
		value[0] = m.Bytes[0]
		value[1] = m.Bytes[1]
		m.Bytes = m.Bytes[2:]
	}
	return nil
}
func (m *mock) WriteToReg(addr, reg byte, value []byte) error { return nil }
func (m *mock) Close() error                                  { return nil }

func mocki2cBus() *mock { return new(mock) }

func TestAds1015Driver(t *testing.T) {
	bus := mocki2cBus() //i2c.MockBus()

	f := Ads1015Factory()
	_, err := f.NewDriver(nil, bus)

	if err == nil {
		t.Error("Adapter creation should fail when configuration is null")
	}

	metadata := f.Metadata()
	if metadata.Name != "ADS1015" {
		t.Error("Incorrect metadata received")
	}

	parameters := f.GetParameters()
	if len(parameters) != 5 {
		t.Error("Incorrect number of parameters received")
	}

	driver, err := f.NewDriver(params, bus)

	if err != nil {
		t.Error(err)
	}

	if driver.Metadata().Name != "ADS1015" {
		t.Error("Unexpected name")
	}

	if !driver.Metadata().HasCapability(hal.AnalogInput) {
		t.Error("analog input cpability should exist")
	}

	if driver.Metadata().HasCapability(hal.DigitalInput) {
		t.Error("Digital input Capability should not exist")
	}

	pins, err := driver.Pins(hal.AnalogInput)
	if err != nil {
		t.Error(err)
	}

	if len(pins) != 4 {
		t.Error("Unexpected number of pins returned by driver")
	}

	pins, err = driver.Pins(hal.DigitalOutput)
	if err == nil {
		t.Error("ADS1015 should not support Digital Output")
	}

	d := driver.(hal.AnalogInputDriver)

	if len(d.AnalogInputPins()) != 4 {
		t.Error("Expected 4 channels")
	}
	if _, err := d.AnalogInputPin(5); err == nil {
		t.Error("Expected error for invalid channel name")
	}

	ch, err := d.AnalogInputPin(0)
	if err != nil {
		t.Error(err)
	}
	if ch.Name() != "0" {
		t.Error("Unexpected channel name")
	}

	_, err = ch.Read()
	if err == nil {
		t.Error("Read should fail due to config mismatch")
	}

	//Set i2c bytes to config and reading 193, 131
	bus.Bytes = []byte{0xC1, 0x83, 0x6F, 0xF0}

	v, err := ch.Read()
	if err != nil {
		t.Error(err)
	}

	if v != 1791 {
		t.Error("Unexepected value")
	}

	if err := d.Close(); err != nil {
		t.Error(err)
	}
}

func TestAds1115Driver(t *testing.T) {
	bus := mocki2cBus() //i2c.MockBus()

	f := Ads1115Factory()
	_, err := f.NewDriver(nil, bus)

	if err == nil {
		t.Error("Adapter creation should fail when configuration is null")
	}

	metadata := f.Metadata()
	if metadata.Name != "ADS1115" {
		t.Error("Incorrect metadata received")
	}

	parameters := f.GetParameters()
	if len(parameters) != 5 {
		t.Error("Incorrect number of parameters received")
	}

	driver, err := f.NewDriver(params, bus)

	if err != nil {
		t.Error(err)
	}

	if driver.Metadata().Name != "ADS1115" {
		t.Error("Unexpected name")
	}

	if !driver.Metadata().HasCapability(hal.AnalogInput) {
		t.Error("analog input cpability should exist")
	}

	if driver.Metadata().HasCapability(hal.DigitalInput) {
		t.Error("Digital input Capability should not exist")
	}

	pins, err := driver.Pins(hal.AnalogInput)
	if err != nil {
		t.Error(err)
	}

	if len(pins) != 4 {
		t.Error("Unexpected number of pins returned by driver")
	}

	pins, err = driver.Pins(hal.DigitalOutput)
	if err == nil {
		t.Error("ADS1015 should not support Digital Output")
	}

	d := driver.(hal.AnalogInputDriver)

	if len(d.AnalogInputPins()) != 4 {
		t.Error("Expected 4 channels")
	}
	if _, err := d.AnalogInputPin(5); err == nil {
		t.Error("Expected error for invalid channel name")
	}

	ch, err := d.AnalogInputPin(0)
	if err != nil {
		t.Error(err)
	}
	if ch.Name() != "0" {
		t.Error("Unexpected channel name")
	}

	_, err = ch.Read()
	if err == nil {
		t.Error("Read should fail due to config mismatch")
	}

	//Set i2c bytes to config and reading 193, 131
	bus.Bytes = []byte{0xC1, 0x83, 0x6F, 0xF0}

	v, err := ch.Read()
	if err != nil {
		t.Error(err)
	}

	if v != 28656 {
		t.Error("Unexepected value")
	}

	if err := d.Close(); err != nil {
		t.Error(err)
	}
}
