package drivers

import (
	"fmt"
	"github.com/reef-pi/rpi/i2c"
	"strconv"
	"strings"
	"time"
)

/*
https://www.atlas-scientific.com/_files/_datasheets/_circuit/pH_EZO_datasheet.pdf
*/

type AtlasEZO struct {
	addr  byte
	bus   i2c.Bus
	delay time.Duration
}

func NewAtlasEZO(addr byte, bus i2c.Bus) *AtlasEZO {
	return &AtlasEZO{
		addr:  addr,
		bus:   bus,
		delay: time.Second,
	}
}

func (a *AtlasEZO) extractIntResponse() (int, error) {
	resp, err := a.read()
	if err != nil {
		return 0, err
	}
	parts := strings.Split(resp, ",")
	if len(parts) != 2 {
		return 0, fmt.Errorf("Malformed response:'%s'", resp)
	}
	return strconv.Atoi(parts[1])
}

func (a *AtlasEZO) extractFloatResponse() (float64, error) {
	resp, err := a.read()
	if err != nil {
		return 0, err
	}
	parts := strings.Split(resp, ",")
	if len(parts) != 2 {
		return 0, fmt.Errorf("Malformed response:'%s'", resp)
	}
	return strconv.ParseFloat(parts[1], 64)
}

func (a *AtlasEZO) command(cmd string) error {
	if err := a.bus.WriteBytes(a.addr, []byte(cmd+"\000")); err != nil {
		return err
	}
	time.Sleep(a.delay)
	return nil
}

func (a *AtlasEZO) read() (string, error) {
	payload, err := a.bus.ReadBytes(a.addr, 31)
	if err != nil {
		return "", err
	}
	if payload[0] != byte(1) {
		return "", fmt.Errorf("Failed to execute. Error:%s", string(payload))
	}
	p := strings.Trim(string(payload[1:]), "\000")
	return p, nil
}

func (a *AtlasEZO) LedOn() error {
	return a.command("L,1")
}

func (a *AtlasEZO) LedOff() error {
	return a.command("L,0")
}

func (a *AtlasEZO) LedState() (bool, error) {
	if err := a.command("L,?"); err != nil {
		return false, err
	}
	i, err := a.extractIntResponse()
	if err != nil {
		return false, err
	}
	return i == 1, nil
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

func (a *AtlasEZO) IsCalibrated() (int, error) {
	if err := a.command("Cal,?"); err != nil {
		return 0, err
	}
	return a.extractIntResponse()
}

func (a *AtlasEZO) Factory() error {
	return a.command("Factory")
}

func (a *AtlasEZO) Find() error {
	return a.command("Find")
}

func (a *AtlasEZO) Information() (string, string, error) {
	if err := a.command("i"); err != nil {
		return "", "", err
	}
	resp, err := a.read()
	if err != nil {
		return "", "", err
	}
	parts := strings.Split(resp, ",")
	if len(parts) != 3 {
		return "", "", fmt.Errorf("Malformed response:%s", resp)
	}
	return parts[1], parts[2], nil
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

func (a *AtlasEZO) Status() (string, string, error) {
	if err := a.command("Status"); err != nil {
		return "", "", err
	}
	//?Status,P,5.038
	resp, err := a.read()
	if err != nil {
		return "", "", err
	}
	parts := strings.Split(resp, ",")
	if len(parts) != 3 {
		return "", "", fmt.Errorf("Malformed response:'%s'", resp)
	}
	return parts[1], parts[2], nil
}

func (a *AtlasEZO) GetTC() (float64, error) {
	if err := a.command("T,?"); err != nil {
		return 0, err
	}
	return a.extractFloatResponse()
}

func (a *AtlasEZO) SetTC(t float64) error {
	return a.command(fmt.Sprintf("T,%f", t))
}
