package cbrapi

import (
	"encoding/xml"
	"strconv"
	"strings"
)

// ValCurs structure for CBR XML parsing
type ValCurs struct {
	XMLName xml.Name   `xml:"ValCurs"` // XML root
	Valute  []Currency `xml:"Valute"`  // list of currencies
}

// Currency structure for parsing each currency with its rate
type Currency struct {
	CharCode  string `xml:"CharCode"`
	VunitRate string `xml:"VunitRate"`
}

// TODO : GetRates(dateReq string) (map[string]float64, error)

func parseRussianFloat(num string) (float64, error) {

	normal := strings.Replace(num, ",", ".", 1)
	result, err := strconv.ParseFloat(normal, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}
