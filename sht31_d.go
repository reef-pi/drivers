package drivers

// https://github.com/hybridgroup/gobot/blob/master/drivers/i2c/sht3x_driver.go
import (
	"github.com/reef-pi/rpi/i2c"
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
}

func NewSHT31D(addr byte, bus i2c.Bus) *sht31d {
	return &sht31d{
		addr:     addr,
		bus:      bus,
		accuracy: SHT3xAccuracyHigh,
		delay:    16 * time.Millisecond,
	}
}

func (s *sht31d) Read() (float64, float64, error) {
	return 0, 0, nil
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
		return fmt.Errof("Invalid accuracy value.")
	}
	s.accuracy = a
	return nil
}

func (s *sht31d) command(cmd []byte, expect int) (read []uint16, err error) {
	read = make([]uint16, expect)

	if _, err = s.bus.WriteBytes(send); err != nil {
		return
	}
	time.Sleep(5 * time.Millisecond)

	buf := make([]byte, 3*expect)
	got, err := s.bus.Read(buf)
	if err != nil {
		return
	}
	if got != (3 * expect) {
		return
	}

	for i := 0; i < expect; i++ {
		crc := crc8.Checksum(buf[i*3:i*3+2], s.crcTable)
		if buf[i*3+2] != crc {
			err = ErrInvalidCrc
			return
		}
		read[i] = uint16(buf[i*3])<<8 | uint16(buf[i*3+1])
	}

	return
}
