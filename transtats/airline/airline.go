package airline

import (
	"encoding/csv"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Airline struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func Download() (<-chan Airline, error) {
	resp, err := http.Get("http://transtats.bts.gov/Download_Lookup.asp?Lookup=L_AIRLINE_ID")
	if err != nil {
		return nil, err
	}
	csvReader := csv.NewReader(resp.Body)
	csvReader.Comma = ','
	csvReader.FieldsPerRecord = 2

	head, err := csvReader.Read()
	if err != nil {
		return nil, err
	}
	if head[0] != "Code" {
		return nil, errors.New("Table Header <Code> is missing!")
	}
	if head[1] != "Description" {
		return nil, errors.New("Table Header <Description> is missing!")
	}

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	out := make(chan Airline)
	go func() {
		for _, entry := range records {
			id, err := strconv.Atoi(entry[0])
			if err != nil {
				log.Fatal(err)
			}
			name := entry[1]
			if i := strings.LastIndex(name, ": "); i > -1 {
				name = name[:i]
			}
			out <- Airline{id, name}
		}
		close(out)
	}()

	return out, nil
}
