package drivers

// https://github.com/hybridgroup/gobot/blob/master/drivers/i2c/sht3x_driver.go
import (
	"fmt"
	"github.com/reef-pi/rpi/i2c"
	"github.com/sigurn/crc8"
	"time"
)

const SHT3xAccuracyLow = 0x16
const SHT3xAccuracyMedium = 0x0b
const SHT3xAccuracyHigh = 0x00

type sht31d struct {
	bus      i2c.Bus
	addr     byte
	accuracy byte
	delay    time.Duration
	crcTable *crc8.Table
}

func NewSHT31D(addr byte, bus i2c.Bus) *sht31d {
	crc8Params := crc8.Params{0x31, 0xff, false, false, 0x00, 0xf7, "CRC-8/SENSIRON"}
	return &sht31d{
		addr:     addr,
		bus:      bus,
		accuracy: SHT3xAccuracyHigh,
		delay:    16 * time.Millisecond,
		crcTable: crc8.MakeTable(crc8Params),
	}
}

func (s *sht31d) Read() (float64, float64, error) {
	ret, err := s.command([]byte{0x24, s.accuracy}, 2)
	if nil != err {
		return 0, 0, err
	}
	// temp  = -49 + 315 * (St / (2^16 - 1))
	t := float64((uint64(3150000)*uint64(ret[0]))/uint64(0xffff)-uint64(490000)) / 10000.0
	// relative humidiy  = 100 * Srh / (2^16 - 1)
	rh := float64((uint64(1000000)*uint64(ret[1]))/uint64(0xffff)) / 10000.0
	return t, rh, nil
}

func (s *sht31d) SerialNumber() (uint32, error) {
	ret, err := s.command([]byte{0x37, 0x80}, 2)
	if err != nil {
		return 0, err
	}
	return (uint32(ret[0]) << 16) | uint32(ret[1]), nil
}

func (s *sht31d) Heater() (bool, error) {
	ret, err := s.command([]byte{0xf3, 0x2d}, 1)
	if err != nil {
		return false, err
	}
	sr := ret[0] // status register
	if err != nil {
		return false, err
	}
	return (1 << 13) == (sr & (1 << 13)), nil
}

func (s *sht31d) SetHeater(enabled bool) error {
	out := []byte{0x30, 0x66}
	if true == enabled {
		out[1] = 0x6d
	}
	return s.bus.WriteBytes(s.addr, out)
}

func (s *sht31d) SetAccuracy(a byte) error {
	switch a {
	case SHT3xAccuracyLow:
		s.delay = 5 * time.Millisecond // Actual max is 4, wait 1 ms longer
	case SHT3xAccuracyMedium:
		s.delay = 7 * time.Millisecond // Actual max is 6, wait 1 ms longer
	case SHT3xAccuracyHigh:
		s.delay = 16 * time.Millisecond // Actual max is 15, wait 1 ms longer
	default:
		return fmt.Errorf("Invalid accuracy value.")
	}
	s.accuracy = a
	return nil
}

func (s *sht31d) command(cmd []byte, expect int) ([]uint16, error) {
	read := make([]uint16, expect)
	if err := s.bus.WriteBytes(s.addr, cmd); err != nil {
		return read, err
	}
	time.Sleep(5 * time.Millisecond)
	buf, err := s.bus.ReadBytes(s.addr, 3*expect)
	if err != nil {
		return read, err
	}
	for i := 0; i < expect; i++ {
		crc := crc8.Checksum(buf[i*3:i*3+2], s.crcTable)
		if buf[i*3+2] != crc {
			return read, fmt.Errorf("Incorrect crc checksum.")
		}
		read[i] = uint16(buf[i*3])<<8 | uint16(buf[i*3+1])
	}
	return read, nil
}
