package drivers

/*
https://www.atlas-scientific.com/_files/_datasheets/_circuit/pH_EZO_datasheet.pdf
*/

import (
	"github.com/reef-pi/rpi/i2c"
)

type AtlasEZO struct {
	addr byte
	bus  i2c.Bus
}

type EZO_STATE byte

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
	return a.bus.WriteBytes(a.addr, payload)
}

func (a *AtlasEZO) Baud() error {
	return nil

}

func (a *AtlasEZO) Calibrate() error {
	return nil
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
	return nil
}

func (a *AtlasEZO) Information() error {
	return nil
}

func (a *AtlasEZO) ChangeI2CAddress() error {
	return nil
}

func (a *AtlasEZO) Led(e EZO_STATE) error {
	return a.command([]byte{byte(e)})
}

func (a *AtlasEZO) ProtocolLock() error {
	return nil
}

func (a *AtlasEZO) Read() error {
	return nil
}

func (a *AtlasEZO) Sleep() error {
	return nil
}

func (a *AtlasEZO) Slope() error {
	return nil
}

func (a *AtlasEZO) Status() error {
	return nil
}

func (a *AtlasEZO) TemperatureCompensate() error {
	return nil
}
