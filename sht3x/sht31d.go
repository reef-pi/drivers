package sht3x

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/reef-pi/rpi/i2c"
	"time"
)

const (
	_delay = 15500 * time.Microsecond
)

var (
	CMD_SINGLE_MEASURE_HIGH = []byte{0x24, 0x00} // Single Measure of Temp. and Hum.; High precise
	CMD_ART                 = []byte{0x2B, 0x32} // Activate "accelerated response time"
	CMD_BREAK               = []byte{0x30, 0x93} // Interrupt "periodic acqusition mode" and return to "single shot mode"
	CMD_RESET               = []byte{0x30, 0xA2} // Soft reset command

)

type payload struct {
	Data [2]byte
	CRC  byte
}

type SHT31D struct {
	addr             byte
	bus              i2c.Bus
	pTemp, pHumidity float64
	pTime            time.Time
}

func (d *SHT31D) read(blockCount int) ([]uint16, error) {
	const blockSize = 2 + 1
	data, err := d.bus.ReadBytes(d.addr, blockCount*blockSize)
	if err != nil {
		return nil, err
	}
	resp := make([]payload, blockCount)
	if err := binary.Read(bytes.NewBuffer(data), binary.BigEndian, resp); err != nil {
		return nil, err
	}

	var result []uint16
	for i := 0; i < blockCount; i++ {
		checksum := crc(0xFF, resp[i].Data[:2])
		if checksum != resp[i].CRC {
			return nil, fmt.Errorf("CRCs doesn't match: CRC from sensor (0x%0X) != calculated CRC (0x%0X)", resp[i].CRC, checksum)
		}
		buf := resp[i].Data[:2]
		v := uint16(buf[0])<<8 + uint16(buf[1])
		result = append(result, v)
	}
	return result, nil
}

func (d *SHT31D) initiateMeasure(cmd []byte) error {
	if err := d.bus.WriteBytes(d.addr, cmd); err != nil {
		return err
	}
	time.Sleep(_delay)
	return nil
}

func (d *SHT31D) ReadSensor() (float64, float64, error) {
	if err := d.initiateMeasure(CMD_SINGLE_MEASURE_HIGH); err != nil {
		return 0, 0, err
	}
	data, err := d.read(2)
	if err != nil {
		return 0, 0, err
	}
	temp := float64(data[0])*175/(0x10000-1) - 45
	rh := float64(data[1]) * 100 / (0x10000 - 1)
	d.pTemp = temp
	d.pHumidity = rh
	d.pTime = time.Now()
	return temp, rh, nil
}

func crc(seed byte, buf []byte) byte {
	for i := 0; i < len(buf); i++ {
		seed ^= buf[i]
		for j := 0; j < 8; j++ {
			if seed&0x80 != 0 {
				seed <<= 1
				seed ^= 0x31
			} else {
				seed <<= 1
			}
		}
	}
	return seed
}

func (s *SHT31D) Temperature() (float64, error) {
	if s.pTime.Before(time.Now().Add(time.Second)) {
		if _, _, err := s.ReadSensor(); err != nil {
			return 0, err
		}
	}
	return s.pTemp, nil
}

func (s *SHT31D) Humidity() (float64, error) {
	if s.pTime.Before(time.Now().Add(time.Second)) {
		if _, _, err := s.ReadSensor(); err != nil {
			return 0, err
		}
	}
	return s.pHumidity, nil
}
