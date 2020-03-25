package mp3

import (
	"fmt"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"github.com/reef-pi/hal"
	"io"
	"os"
	"time"
)

type Driver struct {
	meta hal.Metadata
}

type Config struct {
	Loop bool
}

func NewMP3(c Config) (hal.DigitalOutputDriver, error) {
}

func run(quit chan struct{}) {
	f, err := os.Open("/home/ranjib/Downloads/rr.mp3")
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		fmt.Println("Failed to create decoder. Error:", err)
		return
	}

	p, err := oto.NewPlayer(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		fmt.Println("Failed to create player. Error:", err)
		return
	}
	defer p.Close()
	fmt.Printf("Length: %d[bytes]\n", d.Length())
	buf := make([]byte, 8)
	pos := 0

	for {
		select {
		case <-quit:
			fmt.Println("Quitting. written bytes:", pos)
			return
			return
		default:
			n, err := d.Read(buf)
			if err != nil {
				if err == io.EOF {
					fmt.Println("Last payload length:", n, "written bytes:", pos)
					return
				}
				fmt.Println("Read ERROR:", err)
				return
			}
			if _, err := p.Write(buf); err != nil {
				fmt.Println("Write ERROR:", err)
				return
			}
			pos += n
			if n != 8 {
				fmt.Println("Last payload length:", n, "written bytes:", pos)
				return
			}
		}
	}
}

func main() {
	q := make(chan struct{})
	go run(q)
	time.Sleep(5 * time.Second)
	q <- struct{}{}
}
