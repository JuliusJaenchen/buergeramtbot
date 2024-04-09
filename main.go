package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"

	"github.com/PuerkitoBio/goquery"
)

// This page is the entry link to search for home-regestration
// We need to call it first to get valid cookies
var entryURL = "https://service.berlin.de/terminvereinbarung/termin/all/120686/"

// TODO place current UNIX timestamp in slug. The current one is still hardcoded...
var appointmentURL = "https://service.berlin.de/terminvereinbarung/termin/day/1714428000/"

func main() {
	jar, _ := cookiejar.New(nil)
	httpClient := &http.Client{Jar: jar}

	// make first request to get session-cookies
	req := createRequest(entryURL)
	res, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// make second request to get HTML page with appointments
	req = createRequest(appointmentURL)
	res, err = httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// parse HTML to get open appointments
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	errorMessage := doc.Find(".alert-error")
	if len(errorMessage.Nodes) > 0 {
		log.Fatal("ERROR: Bürgeramt session invalid")
	}

	bookableDataPoints := doc.Find("td.buchbar")
	// TODO parse dates out of links for more details
	if len(bookableDataPoints.Nodes) > 0 {
		fmt.Print("Success! ")
	}
	fmt.Printf("Found %v days with open slots\n", len(bookableDataPoints.Nodes))
}

func createRequest(url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	return req
}
