package tasmota

import (
	"encoding/json"
	"fmt"
	"github.com/reef-pi/hal"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockTasmotaServer creates a mock Tasmota device server for testing
func mockTasmotaServer(t *testing.T) *httptest.Server {
	// In-memory store for device state
	powerStates := make(map[int]bool)
	dimmerStates := make(map[int]float64)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cmnd := r.URL.Query().Get("cmnd")
		if cmnd == "" {
			http.Error(w, "Missing cmnd parameter", http.StatusBadRequest)
			return
		}

		// Parse command - format: "Command<n> value" or "Command value"
		response := make(map[string]interface{})

		// Handle Power commands
		if len(cmnd) >= 5 && cmnd[:5] == "Power" {
			// Extract output number and value
			var outputNum int
			var value string

			// Check if it's Power<n> or just Power
			if len(cmnd) > 5 && cmnd[5] >= '0' && cmnd[5] <= '9' {
				// Power<n> format
				fmt.Sscanf(cmnd[5:], "%d", &outputNum)
				// Find where the number ends
				i := 5
				for i < len(cmnd) && cmnd[i] >= '0' && cmnd[i] <= '9' {
					i++
				}
				if i < len(cmnd) && cmnd[i] == ' ' {
					value = cmnd[i+1:]
				}
			} else if len(cmnd) > 6 && cmnd[5] == ' ' {
				// Power value format
				outputNum = 1
				value = cmnd[6:]
			} else {
				// Query format
				outputNum = 1
				value = ""
			}

			// Handle Power command
			if value == "" {
				// Query
				state, ok := powerStates[outputNum]
				if !ok {
					state = false
				}
				stateStr := "OFF"
				if state {
					stateStr = "ON"
				}
				if outputNum == 0 || outputNum == 1 {
					response["POWER"] = stateStr
				}
				if outputNum != 1 {
					response[fmt.Sprintf("POWER%d", outputNum)] = stateStr
				}
			} else if value == "1" || value == "ON" || value == "on" || value == "true" || value == "True" {
				powerStates[outputNum] = true
				stateStr := "ON"
				if outputNum == 0 || outputNum == 1 {
					response["POWER"] = stateStr
				}
				if outputNum != 1 {
					response[fmt.Sprintf("POWER%d", outputNum)] = stateStr
				}
			} else if value == "0" || value == "OFF" || value == "off" || value == "false" || value == "False" {
				powerStates[outputNum] = false
				stateStr := "OFF"
				if outputNum == 0 || outputNum == 1 {
					response["POWER"] = stateStr
				}
				if outputNum != 1 {
					response[fmt.Sprintf("POWER%d", outputNum)] = stateStr
				}
			}
		} else if len(cmnd) >= 6 && cmnd[:6] == "Dimmer" {
			// Handle Dimmer command
			var value float64
			if len(cmnd) > 7 && cmnd[6] == ' ' {
				fmt.Sscanf(cmnd[7:], "%f", &value)
				dimmerStates[0] = value
				if value > 0 {
					powerStates[0] = true
					response["POWER"] = "ON"
				} else {
					powerStates[0] = false
					response["POWER"] = "OFF"
				}
				response["Dimmer"] = int(value)
			} else {
				// Query
				value, ok := dimmerStates[0]
				if !ok {
					value = 0
				}
				response["Dimmer"] = int(value)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))

	return server
}

func TestHttpDriver_AsDigitalOut(t *testing.T) {

	server := mockTasmotaServer(t)
	defer server.Close()

	// Extract host:port from server URL
	address := server.URL[7:] // Remove "http://"

	f := HttpDriverFactory()

	params := map[string]interface{}{
		"Address": address,
		"Output":  "2",
	}

	d, err := f.NewDriver(params, nil)
	if err != nil {
		t.Fatal(err)
	}

	meta := d.Metadata()
	if len(meta.Capabilities) != 2 {
		t.Error("Expected 2 capabilities, found:", len(meta.Capabilities))
	}

	o, ok := d.(hal.DigitalOutputDriver)
	if !ok {
		t.Error("Failed to type driver to Digital output driver")
	}

	if len(o.DigitalOutputPins()) != 1 {
		t.Error("Expected a single digital output pin, found:", len(o.DigitalOutputPins()))
	}

	p, err := o.DigitalOutputPin(0)
	if err != nil || p == nil {
		t.Error("Expected a digital output pin")
	}

	if p.Name() != "Tasmota" {
		t.Error("Expected Tasmota name, found: ", p.Name())
	}

	if p.Number() != 0 {
		t.Error("Expected number 0, found: ", p.Number())
	}

	// Test with mock server
	err = p.Write(true)
	if err != nil {
		t.Error("Expected write true in the digital output, error: ", err.Error())
	}

	if !p.LastState() {
		t.Error("Expected last state is true")
	}

	err = p.Write(false)
	if err != nil {
		t.Error("Expected write false in the digital output, error: ", err.Error())
	}

	if p.LastState() {
		t.Error("Expected last state is false")
	}

}

func TestHttpDriver_AsPWMDriver(t *testing.T) {

	server := mockTasmotaServer(t)
	defer server.Close()

	// Extract host:port from server URL
	address := server.URL[7:] // Remove "http://"

	f := HttpDriverFactory()

	params := map[string]interface{}{
		"Address": address,
		"Output":  "0",
	}

	d, err := f.NewDriver(params, nil)
	if err != nil {
		t.Fatal(err)
	}

	meta := d.Metadata()
	if len(meta.Capabilities) != 2 {
		t.Error("Expected 2 capabilities, found:", len(meta.Capabilities))
	}

	pwm, ok := d.(hal.PWMDriver)
	if !ok {
		t.Error("Failed to type driver to PWM driver")
	}

	if len(pwm.PWMChannels()) != 1 {
		t.Error("Expected a single pwm channel, found:", len(pwm.PWMChannels()))
	}

	p, err := pwm.PWMChannel(0)
	if err != nil {
		t.Error("Expected a pwm pin")
	}

	if p.Name() != "Tasmota" {
		t.Error("Expected Tasmota name, found: ", p.Name())
	}

	if p.Number() != 0 {
		t.Error("Expected number 0, found: ", p.Number())
	}

	// Test with mock server
	err = p.Set(100)
	if err != nil {
		t.Error("Expected to set 100 in the pwm output, error: ", err.Error())
	}

	if !p.LastState() {
		t.Error("Expected last state is true")
	}

	err = p.Set(0)
	if err != nil {
		t.Error("Expected to set 0 in the pwm output, error: ", err.Error())
	}

	if p.LastState() {
		t.Error("Expected last state is false")
	}

}

func TestParseOutputs_SingleOutput(t *testing.T) {
	outputs, err := parseOutputs("1")
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if len(outputs) != 1 || outputs[0] != 1 {
		t.Errorf("Expected [1], got %v", outputs)
	}
}

func TestParseOutputs_DiscreteOutputs(t *testing.T) {
	outputs, err := parseOutputs("1,2,3")
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if len(outputs) != 3 || outputs[0] != 1 || outputs[1] != 2 || outputs[2] != 3 {
		t.Errorf("Expected [1, 2, 3], got %v", outputs)
	}
}

func TestParseOutputs_Range(t *testing.T) {
	outputs, err := parseOutputs("1-3")
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if len(outputs) != 3 || outputs[0] != 1 || outputs[1] != 2 || outputs[2] != 3 {
		t.Errorf("Expected [1, 2, 3], got %v", outputs)
	}
}

func TestParseOutputs_Mixed(t *testing.T) {
	outputs, err := parseOutputs("1-3,5,7-9")
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	expected := []int{1, 2, 3, 5, 7, 8, 9}
	if len(outputs) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(outputs))
	}
	for i, v := range expected {
		if outputs[i] != v {
			t.Errorf("Expected outputs[%d]=%d, got %d", i, v, outputs[i])
		}
	}
}

func TestParseOutputs_EmptyString(t *testing.T) {
	_, err := parseOutputs("")
	if err == nil {
		t.Error("Expected error for empty string")
	}
}

func TestParseOutputs_Duplicates(t *testing.T) {
	_, err := parseOutputs("1,1")
	if err == nil {
		t.Error("Expected error for duplicate output")
	}
}

func TestParseOutputs_DuplicatesInRange(t *testing.T) {
	_, err := parseOutputs("1-3,2")
	if err == nil {
		t.Error("Expected error for duplicate output in range")
	}
}

func TestParseOutputs_ReversedRange(t *testing.T) {
	_, err := parseOutputs("3-1")
	if err == nil {
		t.Error("Expected error for reversed range")
	}
}

func TestParseOutputs_NegativeNumbers(t *testing.T) {
	_, err := parseOutputs("-1")
	if err == nil {
		t.Error("Expected error for negative number")
	}
}

func TestParseOutputs_InvalidFormat(t *testing.T) {
	_, err := parseOutputs("abc")
	if err == nil {
		t.Error("Expected error for invalid format")
	}
}

func TestParseOutputs_Sorted(t *testing.T) {
	outputs, err := parseOutputs("3,1,2")
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if outputs[0] != 1 || outputs[1] != 2 || outputs[2] != 3 {
		t.Errorf("Expected sorted [1, 2, 3], got %v", outputs)
	}
}

func TestHttpDriver_FactoryValidateParameters(t *testing.T) {

	f := HttpDriverFactory()

	params := map[string]interface{}{
		"Address": "192.168.1.46",
	}

	_, err := f.NewDriver(params, nil)
	if err != nil {
		t.Fatal(err)
	}

	params = map[string]interface{}{
		"Address": "",
	}

	_, err = f.NewDriver(params, nil)
	if err == nil {
		t.Fatal("Expected error")
	}

	params = map[string]interface{}{
		"Address": 1,
	}

	_, err = f.NewDriver(params, nil)
	if err == nil {
		t.Fatal("Expected error")
	}

	params = map[string]interface{}{
		"Address": nil,
	}

	_, err = f.NewDriver(params, nil)
	if err == nil {
		t.Fatal("Expected error")
	}

}

func TestHttpDriver_MultiOutput_DiscreteOutputs(t *testing.T) {
	server := mockTasmotaServer(t)
	defer server.Close()

	address := server.URL[7:] // Remove "http://"

	f := HttpDriverFactory()

	params := map[string]interface{}{
		"Address": address,
		"Output":  "1,2,3",
	}

	d, err := f.NewDriver(params, nil)
	if err != nil {
		t.Fatal("Failed to create driver:", err)
	}

	// Check digital output pins
	dout, ok := d.(hal.DigitalOutputDriver)
	if !ok {
		t.Fatal("Failed to type to DigitalOutputDriver")
	}

	if len(dout.DigitalOutputPins()) != 3 {
		t.Errorf("Expected 3 digital output pins, got %d", len(dout.DigitalOutputPins()))
	}

	// Test each pin
	for i := 0; i < 3; i++ {
		pin, err := dout.DigitalOutputPin(i)
		if err != nil {
			t.Errorf("Failed to get pin %d: %v", i, err)
		}

		// Test Write
		err = pin.Write(true)
		if err != nil {
			t.Errorf("Pin %d: Failed to write true: %v", i, err)
		}

		state := pin.LastState()
		if !state {
			t.Errorf("Pin %d: Expected LastState true, got false", i)
		}

		// Test Write false
		err = pin.Write(false)
		if err != nil {
			t.Errorf("Pin %d: Failed to write false: %v", i, err)
		}

		state = pin.LastState()
		if state {
			t.Errorf("Pin %d: Expected LastState false, got true", i)
		}
	}

	// Check PWM channels
	pwm, ok := d.(hal.PWMDriver)
	if !ok {
		t.Fatal("Failed to type to PWMDriver")
	}

	if len(pwm.PWMChannels()) != 3 {
		t.Errorf("Expected 3 PWM channels, got %d", len(pwm.PWMChannels()))
	}

	// Test each channel - note: Dimmer is a global Tasmota command
	for i := 0; i < 3; i++ {
		ch, err := pwm.PWMChannel(i)
		if err != nil {
			t.Errorf("Failed to get channel %d: %v", i, err)
		}

		// Test Set - verify no errors
		err = ch.Set(100)
		if err != nil {
			t.Errorf("Channel %d: Failed to set 100: %v", i, err)
		}

		// Test Set to 0
		err = ch.Set(0)
		if err != nil {
			t.Errorf("Channel %d: Failed to set 0: %v", i, err)
		}
	}
}

func TestHttpDriver_MultiOutput_Range(t *testing.T) {
	server := mockTasmotaServer(t)
	defer server.Close()

	address := server.URL[7:] // Remove "http://"

	f := HttpDriverFactory()

	params := map[string]interface{}{
		"Address": address,
		"Output":  "1-3",
	}

	d, err := f.NewDriver(params, nil)
	if err != nil {
		t.Fatal("Failed to create driver:", err)
	}

	// Check digital output pins
	dout, ok := d.(hal.DigitalOutputDriver)
	if !ok {
		t.Fatal("Failed to type to DigitalOutputDriver")
	}

	pins := dout.DigitalOutputPins()
	if len(pins) != 3 {
		t.Errorf("Expected 3 digital output pins from range 1-3, got %d", len(pins))
	}

	// Test that each pin can be controlled independently
	for i := 0; i < 3; i++ {
		pin, err := dout.DigitalOutputPin(i)
		if err != nil {
			t.Errorf("Failed to get pin %d: %v", i, err)
		}

		err = pin.Write(true)
		if err != nil {
			t.Errorf("Pin %d: Failed to write: %v", i, err)
		}

		if !pin.LastState() {
			t.Errorf("Pin %d: Expected state true, got false", i)
		}
	}
}

func TestHttpDriver_MultiOutput_OutOfBounds(t *testing.T) {
	server := mockTasmotaServer(t)
	defer server.Close()

	address := server.URL[7:] // Remove "http://"

	f := HttpDriverFactory()

	params := map[string]interface{}{
		"Address": address,
		"Output":  "1,2",
	}

	d, err := f.NewDriver(params, nil)
	if err != nil {
		t.Fatal("Failed to create driver:", err)
	}

	dout, ok := d.(hal.DigitalOutputDriver)
	if !ok {
		t.Fatal("Failed to type to DigitalOutputDriver")
	}

	// Try to access out-of-bounds pin
	_, err = dout.DigitalOutputPin(5)
	if err == nil {
		t.Error("Expected error for out-of-bounds pin access")
	}

	pwm, ok := d.(hal.PWMDriver)
	if !ok {
		t.Fatal("Failed to type to PWMDriver")
	}

	// Try to access out-of-bounds channel
	_, err = pwm.PWMChannel(5)
	if err == nil {
		t.Error("Expected error for out-of-bounds channel access")
	}
}

func TestHttpDriver_BackwardCompatibility_SingleIntegerOutput(t *testing.T) {
	server := mockTasmotaServer(t)
	defer server.Close()

	address := server.URL[7:] // Remove "http://"

	f := HttpDriverFactory()

	// Test with integer output (old format)
	params := map[string]interface{}{
		"Address": address,
		"Output":  1,
	}

	d, err := f.NewDriver(params, nil)
	if err != nil {
		t.Fatal("Failed to create driver with integer output:", err)
	}

	dout, ok := d.(hal.DigitalOutputDriver)
	if !ok {
		t.Fatal("Failed to type to DigitalOutputDriver")
	}

	if len(dout.DigitalOutputPins()) != 1 {
		t.Errorf("Expected 1 pin, got %d", len(dout.DigitalOutputPins()))
	}

	pin, err := dout.DigitalOutputPin(0)
	if err != nil {
		t.Fatal("Failed to get pin:", err)
	}

	err = pin.Write(true)
	if err != nil {
		t.Fatal("Failed to write:", err)
	}

	if !pin.LastState() {
		t.Error("Expected state true")
	}
}

func TestValidation_ValidOutputFormats(t *testing.T) {
	server := mockTasmotaServer(t)
	defer server.Close()

	address := server.URL[7:]
	f := HttpDriverFactory()

	testCases := []struct {
		name   string
		output interface{}
	}{
		{"single digit", 1},
		{"single string", "1"},
		{"discrete outputs", "1,2,3"},
		{"range outputs", "1-5"},
		{"mixed format", "1-3,5,7-9"},
	}

	for _, tc := range testCases {
		params := map[string]interface{}{
			"Address": address,
			"Output":  tc.output,
		}

		_, err := f.NewDriver(params, nil)
		if err != nil {
			t.Errorf("Output format '%v' (%s): unexpected error: %v", tc.output, tc.name, err)
		}
	}
}

func TestValidation_InvalidOutputConfigs(t *testing.T) {
	server := mockTasmotaServer(t)
	defer server.Close()

	address := server.URL[7:]
	f := HttpDriverFactory()

	testCases := []struct {
		name   string
		output interface{}
	}{
		{"empty string", ""},
		{"negative number", -1},
		{"duplicate outputs", "1,1,2"},
		{"reversed range", "5-1"},
		{"invalid format", "abc"},
		{"invalid range", "1-a"},
	}

	for _, tc := range testCases {
		params := map[string]interface{}{
			"Address": address,
			"Output":  tc.output,
		}

		_, err := f.NewDriver(params, nil)
		if err == nil {
			t.Errorf("Output '%v' (%s): expected error but got none", tc.output, tc.name)
		}
	}
}

func TestValidation_MissingAddress(t *testing.T) {
	f := HttpDriverFactory()

	params := map[string]interface{}{
		"Output": "1",
	}

	_, err := f.NewDriver(params, nil)
	if err == nil {
		t.Error("Expected error for missing address")
	}
}

func TestValidation_EmptyAddress(t *testing.T) {
	f := HttpDriverFactory()

	params := map[string]interface{}{
		"Address": "",
		"Output":  "1",
	}

	_, err := f.NewDriver(params, nil)
	if err == nil {
		t.Error("Expected error for empty address")
	}
}

func TestValidation_InvalidAddressType(t *testing.T) {
	f := HttpDriverFactory()

	params := map[string]interface{}{
		"Address": 12345,
		"Output":  "1",
	}

	_, err := f.NewDriver(params, nil)
	if err == nil {
		t.Error("Expected error for non-string address")
	}
}

func TestValidation_MissingOutput(t *testing.T) {
	f := HttpDriverFactory()

	params := map[string]interface{}{
		"Address": "192.168.1.1",
	}

	// Should not error - Output should default to "1"
	d, err := f.NewDriver(params, nil)
	if err != nil {
		t.Errorf("Expected no error with default output, got: %v", err)
	}

	dout, ok := d.(hal.DigitalOutputDriver)
	if !ok {
		t.Fatal("Failed to type to DigitalOutputDriver")
	}

	if len(dout.DigitalOutputPins()) != 1 {
		t.Errorf("Expected 1 pin with default output, got %d", len(dout.DigitalOutputPins()))
	}
}

// mockErrorTasmotaServer creates a mock Tasmota server that returns errors
func mockErrorTasmotaServer(t *testing.T, statusCode int) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte("Error"))
	}))
	return server
}

// mockMalformedJsonServer creates a mock server that returns malformed JSON
func mockMalformedJsonServer(t *testing.T) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{invalid json"))
	}))
	return server
}

func TestEdgeCase_SingleOutput_ZeroIndex(t *testing.T) {
	server := mockTasmotaServer(t)
	defer server.Close()

	address := server.URL[7:]

	f := HttpDriverFactory()
	params := map[string]interface{}{
		"Address": address,
		"Output":  "0",
	}

	d, err := f.NewDriver(params, nil)
	if err != nil {
		t.Fatal("Failed to create driver with output 0:", err)
	}

	dout, ok := d.(hal.DigitalOutputDriver)
	if !ok {
		t.Fatal("Failed to type to DigitalOutputDriver")
	}

	if len(dout.DigitalOutputPins()) != 1 {
		t.Errorf("Expected 1 pin for output 0, got %d", len(dout.DigitalOutputPins()))
	}

	pin, err := dout.DigitalOutputPin(0)
	if err != nil {
		t.Fatal("Failed to get pin:", err)
	}

	err = pin.Write(true)
	if err != nil {
		t.Fatal("Failed to write:", err)
	}
}

func TestEdgeCase_LargeOutputNumber(t *testing.T) {
	server := mockTasmotaServer(t)
	defer server.Close()

	address := server.URL[7:]

	f := HttpDriverFactory()
	params := map[string]interface{}{
		"Address": address,
		"Output":  "32",
	}

	d, err := f.NewDriver(params, nil)
	if err != nil {
		t.Fatal("Failed to create driver with output 32:", err)
	}

	dout, ok := d.(hal.DigitalOutputDriver)
	if !ok {
		t.Fatal("Failed to type to DigitalOutputDriver")
	}

	if len(dout.DigitalOutputPins()) != 1 {
		t.Errorf("Expected 1 pin for output 32, got %d", len(dout.DigitalOutputPins()))
	}
}

func TestEdgeCase_WideOutputRange(t *testing.T) {
	server := mockTasmotaServer(t)
	defer server.Close()

	address := server.URL[7:]

	f := HttpDriverFactory()
	params := map[string]interface{}{
		"Address": address,
		"Output":  "1-10",
	}

	d, err := f.NewDriver(params, nil)
	if err != nil {
		t.Fatal("Failed to create driver with range 1-10:", err)
	}

	dout, ok := d.(hal.DigitalOutputDriver)
	if !ok {
		t.Fatal("Failed to type to DigitalOutputDriver")
	}

	if len(dout.DigitalOutputPins()) != 10 {
		t.Errorf("Expected 10 pins for range 1-10, got %d", len(dout.DigitalOutputPins()))
	}
}

func TestErrorScenario_HTTPError500(t *testing.T) {
	server := mockErrorTasmotaServer(t, 500)
	defer server.Close()

	address := server.URL[7:]

	f := HttpDriverFactory()
	params := map[string]interface{}{
		"Address": address,
		"Output":  "1",
	}

	d, err := f.NewDriver(params, nil)
	if err != nil {
		t.Fatal("Failed to create driver:", err)
	}

	dout, ok := d.(hal.DigitalOutputDriver)
	if !ok {
		t.Fatal("Failed to type to DigitalOutputDriver")
	}

	pin, err := dout.DigitalOutputPin(0)
	if err != nil {
		t.Fatal("Failed to get pin:", err)
	}

	err = pin.Write(true)
	if err == nil {
		t.Error("Expected error on HTTP 500, got nil")
	}
}

func TestErrorScenario_MalformedJSON(t *testing.T) {
	server := mockMalformedJsonServer(t)
	defer server.Close()

	address := server.URL[7:]

	f := HttpDriverFactory()
	params := map[string]interface{}{
		"Address": address,
		"Output":  "1",
	}

	d, err := f.NewDriver(params, nil)
	if err != nil {
		t.Fatal("Failed to create driver:", err)
	}

	dout, ok := d.(hal.DigitalOutputDriver)
	if !ok {
		t.Fatal("Failed to type to DigitalOutputDriver")
	}

	pin, err := dout.DigitalOutputPin(0)
	if err != nil {
		t.Fatal("Failed to get pin:", err)
	}

	// LastState should return false on malformed JSON
	state := pin.LastState()
	if state {
		t.Error("Expected LastState to return false on malformed JSON")
	}
}

func TestEdgeCase_MultipleOutputsConsistentState(t *testing.T) {
	server := mockTasmotaServer(t)
	defer server.Close()

	address := server.URL[7:]

	f := HttpDriverFactory()
	params := map[string]interface{}{
		"Address": address,
		"Output":  "1,2,3",
	}

	d, err := f.NewDriver(params, nil)
	if err != nil {
		t.Fatal("Failed to create driver:", err)
	}

	dout, ok := d.(hal.DigitalOutputDriver)
	if !ok {
		t.Fatal("Failed to type to DigitalOutputDriver")
	}

	// Set all outputs to true
	for i := 0; i < 3; i++ {
		pin, err := dout.DigitalOutputPin(i)
		if err != nil {
			t.Fatalf("Failed to get pin %d: %v", i, err)
		}

		err = pin.Write(true)
		if err != nil {
			t.Fatalf("Pin %d: Failed to write true: %v", i, err)
		}
	}

	// Verify all are true
	for i := 0; i < 3; i++ {
		pin, err := dout.DigitalOutputPin(i)
		if err != nil {
			t.Fatalf("Failed to get pin %d: %v", i, err)
		}

		state := pin.LastState()
		if !state {
			t.Errorf("Pin %d: Expected state true, got false", i)
		}
	}

	// Set all to false
	for i := 0; i < 3; i++ {
		pin, err := dout.DigitalOutputPin(i)
		if err != nil {
			t.Fatalf("Failed to get pin %d: %v", i, err)
		}

		err = pin.Write(false)
		if err != nil {
			t.Fatalf("Pin %d: Failed to write false: %v", i, err)
		}
	}

	// Verify all are false
	for i := 0; i < 3; i++ {
		pin, err := dout.DigitalOutputPin(i)
		if err != nil {
			t.Fatalf("Failed to get pin %d: %v", i, err)
		}

		state := pin.LastState()
		if state {
			t.Errorf("Pin %d: Expected state false, got true", i)
		}
	}
}
