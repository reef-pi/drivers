package ph_board

import (
	"bytes"
	"encoding/binary"
	"github.com/reef-pi/rpi/i2c"
)

const chName = "0"

type channel struct {
	bus  i2c.Bus
	addr byte
}

func (c *channel) Name() string {
	return chName
}

func (c *channel) Read() (float64, error) {
	if err := c.bus.WriteBytes(c.addr, []byte{0x10}); err != nil {
		return -1, err
	}
	buf := make([]byte, 2)
	if err := c.bus.ReadFromReg(c.addr, 0x0, buf); err != nil {
		return -1, err
	}
	var v int16
	if err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &v); err != nil {
		return -1, err
	}
	return float64(v), nil
}
