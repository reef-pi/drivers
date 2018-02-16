package drivers

import (
	"fmt"
	"github.com/reef-pi/rpi/i2c"
)

const (
	REGISTER_DISPLAY_SETUP = 0x80
	REGISTER_SYSTEM_SETUP  = 0x20
	REGISTER_DIMMING       = 0xE0
	BLINKRATE_OFF          = 0x00
	BLINKRATE_2HZ          = 0x01
	BLINKRATE_1HZ          = 0x02
	BLINKRATE_HALFHZ       = 0x03
)

var digits = map[rune]uint16{
	'0': 63,
	'1': 6,
	'2': 219,
	'3': 207,
	'4': 230,
	'5': 237,
	'6': 253,
	'7': 7,
	'8': 255,
	'9': 239,
	'A': 247,
	'B': 4815,
	'C': 57,
	'D': 4623,
	'E': 249,
	'F': 113,
	'G': 189,
	'H': 246,
	'I': 4617,
	'J': 30,
	'K': 9328,
	'L': 56,
	'M': 1334,
	'N': 8502,
	'O': 63,
	'P': 243,
	'Q': 8255,
	'R': 8435,
	'S': 237,
	'T': 4609,
	'U': 62,
	'V': 3120,
	'W': 10294,
	'X': 11520,
	'Y': 5376,
	'Z': 3081,
	' ': 0,
}

type HT16K33 struct {
	buffer []byte
	bus    i2c.Bus
	addr   byte
}

func NewHT16K33(bus i2c.Bus) *HT16K33 {
	return &HT16K33{
		bus:    bus,
		buffer: make([]byte, 16),
		addr:   byte(0x70),
	}
}

func (h *HT16K33) Setup() error {
	if err := h.bus.WriteToReg(h.addr, REGISTER_SYSTEM_SETUP|0x01, []byte{0x00}); err != nil {
		return err
	}
	if err := h.bus.WriteToReg(h.addr, REGISTER_DIMMING|0, []byte{0x00}); err != nil {
		return err
	}
	if err := h.bus.WriteToReg(h.addr, REGISTER_DISPLAY_SETUP|0x01|(BLINKRATE_OFF<<1), []byte{0x00}); err != nil {
		return err
	}
	if bytes, err := h.bus.ReadBytes(h.addr, 16); err != nil {
		return err
	} else {
		h.buffer = bytes
	}

	return h.bus.WriteToReg(h.addr, 0x00, h.buffer)
}

func (h *HT16K33) Blink() error {
	return h.bus.WriteToReg(h.addr, REGISTER_DISPLAY_SETUP|0x01|(BLINKRATE_HALFHZ<<1), []byte{0x00})
}

func (h *HT16K33) Display(word string) error {
	if len(word) != 4 {
		return fmt.Errorf("word length has to be exactly four character")
	}

	for i := 0; i <= 3; i++ {
		item := digits[rune(word[i])]
		h.buffer[i*2], h.buffer[i*2+1] = byte(item), byte(item>>8)
	}
	return h.bus.WriteToReg(h.addr, 0x00, h.buffer)
}
