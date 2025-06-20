package converter

import (
	"errors"
	"github.com/nikomkinds/CurrencyConverter/internal/cbrapi"
	"strconv"
)

type ConversionResult struct {
	From      string  `json:"from"`
	To        string  `json:"to"`
	Amount    float64 `json:"amount"`
	Converted float64 `json:"converted"`
	Rate      float64 `json:"rate"`
}

func Convert(from, to, amountStr, dateReq string) (*ConversionResult, error) {

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return nil, errors.New("invalid amount format")
	}
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	rates, err := cbrapi.GetRates(dateReq)
	if err != nil {
		return nil, err
	}

	fromRate, ok1 := rates[from]
	toRate, ok2 := rates[to]
	if !ok1 {
		return nil, errors.New("currency from is not supported")
	}
	if !ok2 {
		return nil, errors.New("currency to is not supported")
	}

	converted := amount * fromRate / toRate

	return &ConversionResult{
		from,
		to,
		amount,
		converted,
		toRate / fromRate,
	}, nil
}
