package cbrapi

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
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

	// base URL for the last available date
	URL := "http://www.cbr.ru/scripts/XML_daily.asp"
	// add date if provided
	if dateReq != "" {
		if err := isValidDate(dateReq); err != nil {
			return nil, err
		}
		URL += "?date_req=" + dateReq
	}

	// creating a request with headers
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")

	// creating temp client
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	// executing the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	// reading response in case of a mistake
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server returned %d: %s", resp.StatusCode, string(body))
	}

	// reading correct response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading data: %v", err)
	}

	// log.Printf("Raw response:\n%s", string(body))

	// unmarshal XML
	reader := bytes.NewReader(body)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	var valCurs ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return nil, fmt.Errorf("error unmarshalling XML: %v", err)
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

func isValidDate(date string) error {

	re := regexp.MustCompile(`^(0[1-9]|[12][0-9]|3[01])/(0[1-9]|1[0-2])/([0-9]{4})$`)
	if !re.MatchString(date) {
		return errors.New("invalid date format")
	}

	dateYearStr := date[6:]
	dateYear, _ := strconv.Atoi(dateYearStr)
	year := time.Now().Year()
	if dateYear > year {
		return errors.New("invalid year")
	}

	return nil
}

func parseRussianFloat(num string) (float64, error) {

	normal := strings.Replace(num, ",", ".", 1)
	result, err := strconv.ParseFloat(normal, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}
