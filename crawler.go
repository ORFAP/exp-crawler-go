// crawler project crawler.go
package main

import (
	"bytes"
	"github.com/ORFAP/exp-crawler-go/transtats/airline"
	"github.com/ORFAP/exp-crawler-go/transtats/market"
	"github.com/ORFAP/exp-crawler-go/transtats/route"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	postRoutes()
}

func postAirlines() {
	c, _ := airline.Download()
	for a := range c {
		b, _ := json.Marshal(a)
		http.Post("http://10.28.2.166/api/airlines", "application/json", bytes.NewReader(b))
	}
}

func postMarkets() {
	c, _ := market.Download()
	for m := range c {
		b, _ := json.Marshal(m)
		http.Post("http://10.28.2.166/api/markets", "application/json", bytes.NewReader(b))
	}
}

func postRoutes() {
	c, err := route.DownloadForT_ONTIME(2015, 1) //.DownloadForT_ONTIME(2015, 1)

	if err != nil {
		log.Fatal(err)
	}

	var buffer bytes.Buffer
	buffer.WriteRune('[')

	routes := []route.Route{}

	for route := range c {
		routes = append(routes, route)
	}

	b, err := json.Marshal(routes)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", string(b))
	r := bytes.NewReader(b)
	resp, err := http.Post("http://10.28.2.166/api/routes/saveAll", "application/json", r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v", resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("%v", string(body))
}
