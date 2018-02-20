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

func (a *AtlasEZO) command(cmd string) error {
	if err := a.bus.WriteBytes(a.addr, []byte(cmd+"\000")); err != nil {
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
	fmt.Println(payload)
	fmt.Println(payload[0])
	fmt.Println(string(payload[1:]))
	if payload[0] != byte(1) {
		//	return "", fmt.Errorf("Failed to execute. Error:%s", string(payload))
	}
	return string(payload[1:]), nil
}

func (a *AtlasEZO) Baud(n int) error {
	return a.command(fmt.Sprintf("Baud,%f", n))
}

func (a *AtlasEZO) CalibrateMid(n float32) error {
	if err := a.command(fmt.Sprintf("Cal,mid,%f", n)); err != nil {
		return err
	}
	time.Sleep(600 * time.Millisecond)
	return nil
}

func (a *AtlasEZO) CalibrateHigh(n float32) error {
	if err := a.command(fmt.Sprintf("Cal,high,%f", n)); err != nil {
		return err
	}
	time.Sleep(600 * time.Millisecond)
	return nil
}

func (a *AtlasEZO) CalibrateLow(n float32) error {
	if err := a.command(fmt.Sprintf("Cal,low,%f", n)); err != nil {
		return err
	}
	time.Sleep(600 * time.Millisecond)
	return nil
}

func (a *AtlasEZO) ClearCalibration() error {
	return a.command("Cal,clear")
}

func (a *AtlasEZO) IsCalibrated() error {
	return a.command("Cal,?")
}

func (a *AtlasEZO) Export() error {
	return nil
}

func (a *AtlasEZO) Import() error {
	return nil
}

func (a *AtlasEZO) Factory() error {
	return a.command("Factory")
}

func (a *AtlasEZO) Find() error {
	return a.command("Find")
}

func (a *AtlasEZO) Information() error {
	return a.command("i")
}

func (a *AtlasEZO) ChangeI2CAddress() error {
	return nil
}

func (a *AtlasEZO) SetLed(e EZO_STATE) error {
	return a.command(string([]byte{byte(e)}))
}

func (a *AtlasEZO) GetLed() error {
	return a.command("L,?")
}

func (a *AtlasEZO) ProtocolLock() error {
	return nil
}

func (a *AtlasEZO) Read() (float64, error) {
	if err := a.command("R"); err != nil {
		return 0, err
	}
	v, err := a.read()
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(v, 64)
}

func (a *AtlasEZO) Sleep() error {
	return a.command("Sleep")
}

func (a *AtlasEZO) Slope() error {
	return nil
}

func (a *AtlasEZO) Status() error {
	return a.command("Status")
}

func (a *AtlasEZO) TemperatureCompensate() error {
	return nil
}
