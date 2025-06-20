package converter

type ConversionResult struct {
	From      string  `json:"from"`
	To        string  `json:"to"`
	Amount    float64 `json:"amount"`
	Converted float64 `json:"converted"`
	Rate      float64 `json:"rate"`
}

// TODO : Convert (from, to, amount string) (*ConversionResult, error)
