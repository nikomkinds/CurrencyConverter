package cbrapi

import (
	"encoding/xml"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// ValCurs structure for CBR XML parsing
type ValCurs struct {
	XMLName xml.Name   `xml:"ValCurs"` // XML root
	Valutes []Currency `xml:"Valute"`  // list of currencies
}

// Currency structure for parsing each currency with its rate
type Currency struct {
	CharCode  string `xml:"CharCode"`
	VunitRate string `xml:"VunitRate"`
}

func GetRates(dateReq string) (map[string]float64, error) {

	// base url for the last available date
	url := "http://www.cbr.ru/scripts/XML_daily.asp"
	// add date if provided
	if dateReq != "" && isValidDate(dateReq) {
		url += "?date_req=" + dateReq
	}

	// getting data from CBR
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// unmarshal XML
	var valCurs ValCurs
	if err := xml.Unmarshal(body, &valCurs); err != nil {
		return nil, err
	}

	// reading all currencies
	rates := make(map[string]float64)
	rates["RUB"] = 1.0
	for _, v := range valCurs.Valutes {

		value, err := parseRussianFloat(v.VunitRate)
		if err != nil {
			continue
		}
		rates[v.CharCode] = value
	}
	return rates, nil
}

func isValidDate(date string) bool {
	re := regexp.MustCompile(`^(0[1-9]|[12][0-9]|3[01])/(0[1-9]|1[0-2])/([0-9]{4})&`)
	return re.MatchString(date)
}

func parseRussianFloat(num string) (float64, error) {

	normal := strings.Replace(num, ",", ".", 1)
	result, err := strconv.ParseFloat(normal, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}
