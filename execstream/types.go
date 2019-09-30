package execstream

type Command struct {
	Type string `json:"type"`
	Pin int `json:"pin"`
	Value *float64 `json:"value,omitempty"`
}

type Response struct {
	Type string `json:"type"`
	Pin int `json:"pin"`
	Value *float64 `json:"value,omitempty"`
}
