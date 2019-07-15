package tplink

type (
	Action struct {
		Type float64 `json:"type,omitempty"`
	}

	Child struct {
		ID         string `json:"id,omitempty"`
		State      int    `json:"state,omitempty"`
		Alias      string `json:"alias,omitempty"`
		OnTime     int    `json:"on_time,omitempty"`
		NextAction Action `json:"next_action,omitempty"`
	}

	Sysinfo struct {
		Alias           string  `json:"alias,omitempty"`
		SoftwareVersion string  `json:"sw_veri,omitempty"`
		HardwareVersion string  `json:"hw_ver,omitempty"`
		Model           string  `json:"model,omitempty"`
		DeviceID        string  `json:"deviceId,omitempty"`
		OemID           string  `json:"oemId,omitempty"`
		HardwareID      string  `json:"hwId,omitempty"`
		Rssi            float64 `json:"rssi,omitempty"`
		Longitude       float64 `json:"longitude,omitempty"`
		Latitude        float64 `json:"latitude,omitempty"`
		Updating        int     `json:"updating,omitempty"`
		LEDOff          int     `json:"led_off,omitempty"`
		RelayState      int     `json:"relay_state,omitempty"`
		OnTime          int     `json:"on_time,omitempty"`
		ActiveMode      string  `json:"active_mode,omitempty"`
		IconHash        string  `json:"icon_hash,omitempty"`
		ErrorCode       int     `json:"err_code,omitempty"`
		Children        []Child `json:"children,omitempty"`
	}

	System struct {
		Sysinfo Sysinfo `json:"get_sysinfo"`
	}
	Plug struct {
		System System `json:"system"`
	}
	Config struct {
		Address string `json:"address"`
	}
	CmdRelayState struct {
		System struct {
			RelayState struct {
				State int `json:"state"`
			} `json:"set_relay_state"`
		} `json:"system"`
		Context struct {
			Children []string `json:"child_ids,omitempty"`
		} `json:"context,omitempty"`
	}
)
