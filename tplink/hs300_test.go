package tplink

import (
	"testing"
)

func TestHS300Strip(t *testing.T) {
	d := NewHS300Strip("127.0.0.1:9999")
	d.cnFactory = mockConnFacctory
	if d.Metadata().Name == "" {
		t.Error("HAL metadata should not have empty name")
	}

}
