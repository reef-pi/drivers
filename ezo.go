package drivers

import (
	"fmt"
	"github.com/reef-pi/rpi/i2c"
	"strconv"
	"time"
)

/*
https://www.atlas-scientific.com/_files/_datasheets/_circuit/pH_EZO_datasheet.pdf
*/

type AtlasEZO struct {
	addr byte
	bus  i2c.Bus
}

type EZO_STATE byte

const EZO_PAUSE_TIME = 300 * time.Millisecond

const (
	EZO_OFF EZO_STATE = iota
	EZO_ON
)

func NewAtlasEZO(addr byte, bus i2c.Bus) *AtlasEZO {
	return &AtlasEZO{
		addr: addr,
		bus:  bus,
	}
}

func (a *AtlasEZO) command(payload []byte) error {
	if err := a.bus.WriteBytes(a.addr, payload); err != nil {
		return err
	}
	time.Sleep(EZO_PAUSE_TIME)
	return nil
}

func (a *AtlasEZO) read() (string, error) {
	payload, err := a.bus.ReadBytes(a.addr, 31)
	if err != nil {
		return "", err
	}
	if payload[0] != byte('1') {
		return "", fmt.Errorf("Failed to execute. Error:%s", string(payload))
	}
	return string(payload[1:]), nil
}

func (a *AtlasEZO) Baud() error {
	return nil
}

func (a *AtlasEZO) CalibrateMid(n float32) error {
	if err := a.command([]byte(fmt.Sprintf("Cal,mid,%f", n))); err != nil {
		return err
	}
	time.Sleep(600 * time.Millisecond)
	return nil
}

func (a *AtlasEZO) CalibrateHigh(n float32) error {
	if err := a.command([]byte(fmt.Sprintf("Cal,high,%f", n))); err != nil {
		return err
	}
	time.Sleep(600 * time.Millisecond)
	return nil
}

func (a *AtlasEZO) CalibrateLow(n float32) error {
	if err := a.command([]byte(fmt.Sprintf("Cal,low,%f", n))); err != nil {
		return err
	}
	time.Sleep(600 * time.Millisecond)
	return nil
}

func (a *AtlasEZO) ClearCalibration() error {
	return a.command([]byte("Cal,clear"))
}

func (a *AtlasEZO) IsCalibrated() error {
	return a.command([]byte("Cal,?"))
}

func (a *AtlasEZO) Export() error {
	return nil
}

func (a *AtlasEZO) Import() error {
	return nil
}

func (a *AtlasEZO) Factory() error {
	return nil
}

func (a *AtlasEZO) Find() error {
	return a.command([]byte("Find"))
}

func (a *AtlasEZO) Information() error {
	return a.command([]byte("i"))
}

func (a *AtlasEZO) ChangeI2CAddress() error {
	return nil
}

func (a *AtlasEZO) SetLed(e EZO_STATE) error {
	return a.command([]byte{byte(e)})
}

func (a *AtlasEZO) GetLed() error {
	return a.command([]byte("?"))
}

func (a *AtlasEZO) ProtocolLock() error {
	return nil
}

func (a *AtlasEZO) Read() (float64, error) {
	if err := a.command([]byte("R")); err != nil {
		return 0, err
	}
	v, err := a.read()
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(v, 64)
}

func (a *AtlasEZO) Sleep() error {
	return nil
}

func (a *AtlasEZO) Slope() error {
	return nil
}

func (a *AtlasEZO) Status() error {
	return a.command([]byte("Status"))
}

func (a *AtlasEZO) TemperatureCompensate() error {
	return nil
}
