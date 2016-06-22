package route

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Route struct {
	Date           string  `json:"date"`
	Delays         float64 `json:"delays"`
	Cancelled      float64 `json:"cancelled"`
	PassengerCount float64 `json:"passengerCount"`
	FlightCount    float64 `json:"flightCount"`
	Airline        string  `json:"airline"`
	Source         string  `json:"source"`
	Destination    string  `json:"destination"`
}

func DownloadForT_ONTIME(year int, month int) (<-chan Route, error) {

	sql := fmt.Sprintf(`SELECT FL_DATE,ARR_DELAY_NEW,CANCELLED,AIRLINE_ID,ORIGIN_CITY_MARKET_ID,DEST_CITY_MARKET_ID
						FROM T_ONTIME 
						WHERE Month=%v 
						AND YEAR=%v 
						AND ORIGIN_CITY_MARKET_ID=31703`, strconv.Itoa(month), strconv.Itoa(year))

	resp, err := http.PostForm("http://transtats.bts.gov/DownLoad_Table.asp", url.Values{"sqlstr": {sql}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	readerAt := bytes.NewReader(b)
	routesZip, err := zip.NewReader(readerAt, resp.ContentLength)
	if err != nil {
		return nil, err
	}

	csvString, _ := routesZip.File[0].Open()

	csvReader := csv.NewReader(csvString)
	csvReader.Comma = ','

	_, err = csvReader.Read()
	if err != nil {
		return nil, err
	}

	out := make(chan Route)
	go func() {
		for {
			entry, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			delays, _ := strconv.ParseFloat(entry[1], 64)
			cancelled, _ := strconv.ParseFloat(entry[2], 64)
			out <- Route{
				Date:        entry[0],
				Delays:      delays,
				Cancelled:   cancelled,
				Airline:     entry[3],
				Source:      entry[4],
				Destination: entry[5]}
		}
		close(out)
	}()

	return out, nil
}

func DownloadForT_T100D(year int) (<-chan Route, error) {

	sql := fmt.Sprintf(`SELECT MONTH,PASSENGERS,DEPARTURES_PERFORMED,AIRLINE_ID,ORIGIN_CITY_MARKET_ID,DEST_CITY_MARKET_ID
						FROM  T_T100D_SEGMENT_ALL_CARRIER
						WHERE YEAR=%v
						AND ORIGIN_CITY_MARKET_ID=31703`, strconv.Itoa(year))

	resp, err := http.PostForm("http://transtats.bts.gov/DownLoad_Table.asp", url.Values{"sqlstr": {sql}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	readerAt := bytes.NewReader(b)
	routesZip, err := zip.NewReader(readerAt, resp.ContentLength)
	if err != nil {
		return nil, err
	}

	csvString, _ := routesZip.File[0].Open()

	csvReader := csv.NewReader(csvString)
	csvReader.Comma = ','

	_, err = csvReader.Read()
	if err != nil {
		return nil, err
	}

	out := make(chan Route)
	go func() {
		for {
			entry, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			date := fmt.Sprintf("%v-%v-01", year, entry[0])
			passengerCount, _ := strconv.ParseFloat(entry[1], 64)
			flightCount, _ := strconv.ParseFloat(entry[2], 64)
			out <- Route{
				Date:           date,
				PassengerCount: passengerCount,
				FlightCount:    flightCount,
				Airline:        "http://asdf:8080/airlines/" + entry[3],
				Source:         "/markets/" + entry[4],
				Destination:    "/markets/" + entry[5]}
		}
		close(out)
	}()

	return out, nil
}
