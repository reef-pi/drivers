package tplink

import (
	"github.com/reef-pi/hal"
	"testing"
	"time"
)

type mockConn struct{}

func (c *mockConn) Close() error                  { return nil }
func (c *mockConn) SetDeadline(_ time.Time) error { return nil }
func (c *mockConn) Write(_ []byte) (int, error)   { return 0, nil }
func mockConnFacctory(_, _ string, _ time.Duration) (Conn, error) {
	return &mockConn{}, nil
}

func TestHS1xxPlug(t *testing.T) {
	p := NewHS1xxPlug("127.0.0.1:9999")
	p.cnFactory = mockConnFacctory
	if err := p.On(); err != nil {
		t.Error(err)
	}
	if err := p.Off(); err != nil {
		t.Error(err)
	}

	d, err := HALAdapter([]byte(`{"address":"127.0.0.1:3000"}`), nil)
	if err != nil {
		t.Error(err)
	}
	if d.Metadata().Name == "" {
		t.Error("HAL metadata should not have empty name")
	}

	d1, ok := d.(hal.OutputDriver)
	if !ok {
		t.Fatal("Failed to type cast to output driver")
	}

	if len(d1.OutputPins()) != 1 {
		t.Error("Expected exactly one output pin")
	}
	pin, err := d1.OutputPin(0)
	if err != nil {
		t.Error(err)
	}
	if pin.LastState() != false {
		t.Error("Expected initial state to be false")
	}
}
