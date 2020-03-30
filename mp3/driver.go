package mp3

import (
	"fmt"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"github.com/reef-pi/hal"
	"io"
	"log"
	"os"
	"sync"
)

const (
	_name = "mp3"
)

var ctx *oto.Context

type Driver struct {
	sync.Mutex
	quitCh chan struct{}
	state  bool
	meta   hal.Metadata
	conf   Config
}

type Config struct {
	Loop bool
	File string
}

func NewMP3(c Config) (hal.DigitalOutputDriver, error) {
	once.Do(func() {
		context, err := oto.NewContext(44100, 2, 2, 8192)
		if err == nil {
			ctx = context
			return
		}
		log.Println("ERROR: failed to initialize mp3 player context.", err)
	})
	if ctx == nil {
		return nil, fmt.Errorf("mp3 player context not initialized")
	}
	return &Driver{
		meta: hal.Metadata{
			Name: _name,
		},
		conf: c,
	}, nil
}

func (d *Driver) run(quit chan struct{}) {
	f, err := os.Open(d.conf.File)
	if err != nil {
		log.Println("ERROR: failed to open mp3 file", err)
		return
	}
	defer f.Close()

	dec, err := mp3.NewDecoder(f)
	if err != nil {
		log.Println("ERROR: Failed to create mp3 decoder:", err)
		return
	}

	p := ctx.NewPlayer()
	defer p.Close()
	buf := make([]byte, 8)
	for {
		select {
		case <-quit:
			return
		default:
			_, err := dec.Read(buf)
			if err != nil {
				if err == io.EOF {
					if !d.conf.Loop {
						return
					}
					f.Seek(0, io.SeekStart)
					dec, err = mp3.NewDecoder(f)
					if err != nil {
						log.Println("ERRPR: failed to recreate mp3 decoder:", err)
						return
					}
				}
				log.Println("ERROR: mp3 decoder read failed:", err)
				return
			}
			if _, err := p.Write(buf); err != nil {
				log.Println("ERROR: mp3 player write failed:", err)
				return
			}
		}
	}
}

func (d *Driver) Metadata() hal.Metadata {
	return d.meta
}

func (d *Driver) Name() string {
	return d.meta.Name
}

func (d *Driver) Number() int {
	return 0
}
func (d *Driver) DigitalOutputPins() []hal.DigitalOutputPin {
	return []hal.DigitalOutputPin{d}
}

func (d *Driver) DigitalOutputPin(i int) (hal.DigitalOutputPin, error) {
	if i != 0 {
		return nil, fmt.Errorf("invalid pin: %d", i)
	}
	return d, nil
}

func (d *Driver) Write(state bool) error {
	if state {
		return d.On()
	}
	return d.Off()
}
func (d *Driver) On() error {
	if ctx == nil {
		return fmt.Errorf("mp3 player context not initialized")
	}
	d.Lock()
	defer d.Unlock()

	if d.quitCh != nil {
		return fmt.Errorf("previous invoke is still running")
	}
	d.quitCh = make(chan struct{})
	go d.run(d.quitCh)
	d.state = true
	return nil
}
func (d *Driver) Off() error {
	d.state = false
	if d.quitCh == nil {
		return nil
	}
	d.quitCh <- struct{}{}
	close(d.quitCh)
	d.quitCh = nil
	return nil
}

func (d *Driver) LastState() bool {
	return d.state
}

func (d *Driver) Close() error {
	return nil
}
func (d *Driver) Pins(cap hal.Capability) ([]hal.Pin, error) {
	switch cap {
	case hal.DigitalOutput:
		return []hal.Pin{d}, nil
	default:
		return nil, fmt.Errorf("unsupported capability:%s", cap.String())
	}
}
