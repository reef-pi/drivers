package ph_board

import (
	"encoding/binary"
	"github.com/reef-pi/rpi/i2c"
	"math"
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
	buf, err := c.bus.ReadBytes(c.addr, 2)
	if err != nil {
		return -1, err
	}
	bits := binary.LittleEndian.Uint16(buf)
	return math.Float64frombits(uint64(bits)), nil
}
