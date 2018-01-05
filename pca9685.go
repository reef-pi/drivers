package drivers

import (
	"github.com/reef-pi/rpi/i2c"
	"math"
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

func NewPCA9685(addr byte, bus i2c.Bus) *PCA9685 {
	return &PCA9685{
		addr: addr,
		bus:  bus,
		Freq: defaultFreq,
	}
}

func (p *PCA9685) i2cWrite(reg byte, payload []byte) error {
	return p.bus.WriteToReg(p.addr, reg, payload)
}

func (p *PCA9685) mode1Reg() (byte, error) {
	mode1Reg := make([]byte, 1)
	return mode1Reg[0], p.i2cWrite(mode1RegAddr, mode1Reg)
}

func (p *PCA9685) sleep() error {
	mode1Reg, err := p.mode1Reg()
	if err != nil {
		return err
	}
	sleepmode := (mode1Reg & 0x7F) | 0x10
	return p.bus.WriteToReg(p.addr, mode1RegAddr, []byte{sleepmode})
}

func (p *PCA9685) Wake() error {
	mode1Reg, err := p.mode1Reg()
	if err != nil {
		return err
	}
	if err := p.sleep(); err != nil {
		return err
	}
	if p.Freq == 0 {
		p.Freq = defaultFreq
	}
	preScaleValue := byte(math.Floor(float64(clockFreq/(pwmControlPoints*p.Freq))+float64(0.5)) - 1)
	if err := p.bus.WriteToReg(p.addr, preScaleRegAddr, []byte{preScaleValue}); err != nil {
		return err
	}
	newmode := ((mode1Reg | 0x01) & 0xDF)
	return p.bus.WriteToReg(p.addr, mode1RegAddr, []byte{newmode})
}

func (p *PCA9685) SetPwm(channel, onTime, offTime int) error {
	onTimeLowReg := byte(pwm0OnLowReg + (4 * channel))
	onTimeLow := byte(onTime & 0xFF)
	onTimeHigh := byte(onTime >> 8)
	offTimeLow := byte(offTime & 0xFF)
	offTimeHigh := byte(offTime >> 8)
	if err := p.bus.WriteToReg(p.addr, onTimeLowReg, []byte{onTimeLow}); err != nil {
		return err
	}
	onTimeHighReg := onTimeLowReg + 1
	if err := p.bus.WriteToReg(p.addr, onTimeHighReg, []byte{onTimeHigh}); err != nil {
		return err
	}

	offTimeLowReg := onTimeHighReg + 1
	if err := p.bus.WriteToReg(p.addr, offTimeLowReg, []byte{offTimeLow}); err != nil {
		return err
	}

	offTimeHighReg := offTimeLowReg + 1
	return p.bus.WriteToReg(p.addr, offTimeHighReg, []byte{offTimeHigh})
}

func (p *PCA9685) Close() error {
	if err := p.bus.WriteToReg(p.addr, mode1RegAddr, []byte{0x00}); err != nil {
		return err
	}
	for regAddr := 0x06; regAddr <= 0x45; regAddr++ {
		if err := p.bus.WriteToReg(p.addr, byte(regAddr), []byte{0x00}); err != nil {
			return err
		}
	}
	return nil
}
