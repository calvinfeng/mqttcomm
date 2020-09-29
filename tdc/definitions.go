package tdc

type MQTTClientOptions struct {
	BrokerAddr string
}

type MQTTClient interface {
	Connect() error
	SubmitDIOState(state DIOState) error
}

type (
	// DIOState is known as the digital input/output state of a GPIO pin.
	DIOState struct {
		PinName             string   `json:"DioName"`
		Value               float64  `json:"Value"`
		Direction           string   `json:"Direction"`
		AvailableDirections []string `json:"AvailableDirections"`
	}

	LEDState struct {
		ID    string `json:"Id"`
		State string `json:"State"`
	}

	AINState struct {
		PinName string  `json:"AinName"`
		Value   float64 `json:"Value"`
		Mode    string  `json:"Mode"`
		State   string  `json:"State"`
	}

	AINModeChange struct {
		PinName      string `json:"AinName"`
		PreviousMode string `json:"PreviousMode"`
		NewMode      string `json:"NewMode"`
	}

	AINStateChange struct {
		PinName       string `json:"AinName"`
		PreviousState string `json:"PreviousState"`
		NewState      string `json:"NewState"`
	}

	AINValueChange struct {
		PinName       string  `json:"AinName"`
		PreviousValue float64 `json:"PreviousValue"`
		NewValue      float64 `json:"NewValue"`
	}

	DIODirectionChange struct {
		PinName   string `json:"AinName"`
		Direction string `json:"Direction"`
	}
)
