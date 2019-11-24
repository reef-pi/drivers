package pca9685

import (
	"log"
	"math"
	"time"

	"github.com/reef-pi/rpi/i2c"
)

const (
	clockFreq        = 25000000
	pwmControlPoints = 4096
	mode1RegAddr     = 0x00
	preScaleRegAddr  = 0xFE
	pwm0OnLowReg     = 0x6
	defaultFreq      = 490
)

type PCA9685 struct {
	addr byte
	bus  i2c.Bus
	Freq int
}

func New(addr byte, bus i2c.Bus) *PCA9685 {
	return &PCA9685{
		addr: addr,
		bus:  bus,
		Freq: defaultFreq,
	}
}

func (p *PCA9685) mode1Reg() (byte, error) {
	mode1Reg := make([]byte, 1)
	return mode1Reg[0], p.bus.WriteToReg(p.addr, mode1RegAddr, mode1Reg)
}

// Set the sleep flag on the PCA. This will shut down the oscillators.
func (p *PCA9685) Sleep() error {
	mode1Reg, err := p.mode1Reg()
	if err != nil {
		return err
	}

	sleepmode := (mode1Reg & 0x7F) | 0x10 // Mask restart bit and set sleep bit
	return p.bus.WriteToReg(p.addr, mode1RegAddr, []byte{sleepmode})
}

func (p *PCA9685) Wake() error {
	mode1Reg, err := p.mode1Reg()
	if err != nil {
		return err
	}
	if err := p.Sleep(); err != nil {
		return err
	}
	if p.Freq == 0 {
		p.Freq = defaultFreq
	}
	preScaleValue := byte(math.Floor(float64(clockFreq/(pwmControlPoints*p.Freq))+float64(0.5)) - 1)
	if err := p.bus.WriteToReg(p.addr, preScaleRegAddr, []byte{preScaleValue}); err != nil {
		return err
	}
	wakeMode := mode1Reg & 0xEF
	if (mode1Reg & 0x80) == 0x80 {
		if err := p.bus.WriteToReg(p.addr, mode1RegAddr, []byte{wakeMode}); err != nil {
			return err
		}
		time.Sleep(500 * time.Microsecond)
	}

	restartOpCode := wakeMode | 0x80
	if err := p.bus.WriteToReg(p.addr, mode1RegAddr, []byte{restartOpCode}); err != nil {
		return err
	}

	newmode := ((mode1Reg | 0x01) & 0xDF)
	return p.bus.WriteToReg(p.addr, mode1RegAddr, []byte{newmode})
}

func (p *PCA9685) SetPwm(channel int, onTime, offTime uint16) error {
	log.Println("onTime ", onTime, " offTime ", offTime)
	// Split the ints into 4 bytes
	timeReg := byte(pwm0OnLowReg + (4 * channel))
	onTimeLow := byte(onTime & 0xFF)
	onTimeHigh := byte(onTime >> 8)
	offTimeLow := byte(offTime & 0xFF)
	offTimeHigh := byte(offTime >> 8)

	//log.Println("onLow ", onTimeLow, " onHigh ", onTimeHigh, " offLow ", offTimeLow, " offHigh ", offTimeHigh)
	if err := p.bus.WriteToReg(p.addr, timeReg, []byte{onTimeLow}); err != nil {
		return err
	}
	if err := p.bus.WriteToReg(p.addr, timeReg+1, []byte{onTimeHigh}); err != nil {
		return err
	}
	if err := p.bus.WriteToReg(p.addr, timeReg+2, []byte{offTimeLow}); err != nil {
		return err
	}
	return p.bus.WriteToReg(p.addr, timeReg+3, []byte{offTimeHigh})
}

func (p *PCA9685) Close() error {
	// Clear all channels to full off
	for regAddr := 0x06; regAddr < 0x50; regAddr += 4 {
		if err := p.bus.WriteToReg(p.addr, byte(regAddr), []byte{0x00, 0x00, 0x00, 0x10}); err != nil {
			return err
		}
	}
	return nil
}
