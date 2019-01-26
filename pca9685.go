package drivers

import (
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

func NewPCA9685(addr byte, bus i2c.Bus) *PCA9685 {
	return &PCA9685{
		addr: addr,
		bus:  bus,
		Freq: defaultFreq,
	}
}

// Commented because it's not currently used
//func (p *PCA9685) i2cWrite(reg byte, payload []byte) error {
//	return p.bus.WriteToReg(p.addr, reg, payload)
//}

func (p *PCA9685) i2cRead(reg byte, payload []byte) error {
	return p.bus.ReadFromReg(p.addr, reg, payload)
}

func (p *PCA9685) mode1Reg() (byte, error) {
	mode1Reg := make([]byte, 1)
	return mode1Reg[0], p.i2cRead(mode1RegAddr, mode1Reg)
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
	// Read mode1 register
	mode1Reg, err := p.mode1Reg()
	if err != nil {
		return err
	}
	
	if (mode1Reg & 0x80) != 0 {
		// We are in sleep mode after a previous run without shutdown. Restore.
		// First, clear sleep bit
		mode1Reg &= (^byte(0x10))
		p.bus.WriteToReg(p.addr, mode1RegAddr, []byte{mode1Reg})
		// Allow oscillator to stabilize
		time.Sleep(500 * time.Microsecond)
		// Clear sleep bit
		p.bus.WriteToReg(p.addr, mode1RegAddr, []byte{mode1Reg | 0x80})
	} else if (mode1Reg & 0x10) != 0 {
		// We are in normal sleep, do a normal wakeup
		mode1Reg &= (^byte(0x10))
		p.bus.WriteToReg(p.addr, mode1RegAddr, []byte{mode1Reg})
		// Allow oscillator to stabilize
		time.Sleep(500 * time.Microsecond)
	}
	if p.Freq == 0 {
		p.Freq = defaultFreq
	}
	preScaleValue := byte(math.Floor(float64(clockFreq/(pwmControlPoints*p.Freq))+float64(0.5)) - 1)
	if err := p.bus.WriteToReg(p.addr, preScaleRegAddr, []byte{preScaleValue}); err != nil {
		return err
	}
	
	// Set our operating modes:
	mode1Reg = 0x20 // No AllCall, no subaddresses, no sleep, internal clock, enable auto increment
	
	return p.bus.WriteToReg(p.addr, mode1RegAddr, []byte{mode1Reg})
}

func (p *PCA9685) SetPwm(channel, onTime, offTime int) error {
	// At this pont onTime and offTime are alreeady scaled to 0 .. 4096 by the HAL.
	// The PCA9685 has two special states, full on and full off, besides the normal PWM.
	// Using them prevents the microspikes that can cause extra heat generation in mosfet
	// output stages as well as switching noise.
	// Generally, if onTime + 1 == offTime, we're dealing with full on. If onTime == offTome, 
	// it's full off.
	// Since onTime is 0 and always be 0, and offTime will vary between 0 .. 4096, which is
	// one step out of range, we can use that as an indicator.
	// 100 * 40.96 will result in 4096. This triggers a potential issue because LEDx_OFF_H(4)
	// Is the full off flag bit, making 4096 (0x1000) result in full off!
	// Because of that, we need to clamp it here
	if offTime >= 4095 {
		offTime = 4095;
	} 
	
	// If offTime == 0, we want to be full off. Set LEDx_OFF_H(4)
	if offTime == 0 {
		offTime = 4096
	}
	
	// If offTime == 4095, we want to be full on. Set LEDx_ON_H(4)
	if offTime == 4095 {
		onTime = 4096
	} 

	// Split the ints into 4 bytes	
	timeReg := byte(pwm0OnLowReg + (4 * channel))
	onTimeLow := byte(onTime & 0xFF)
	onTimeHigh := byte(onTime >> 8)
	offTimeLow := byte(offTime & 0xFF)
	offTimeHigh := byte(offTime >> 8)
	
	// Send one entire channel in one go
	return p.bus.WriteToReg(p.addr, timeReg, []byte{onTimeLow, onTimeHigh, offTimeLow, offTimeHigh})
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
