package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/joho/godotenv/autoload"
)

// This page is the entry link to search for home-regestration
// We need to call it first to get valid cookies
var entryURL = "https://service.berlin.de/terminvereinbarung/termin/all/120686/"

// TODO place current UNIX timestamp in slug. The current one is still hardcoded...
var appointmentURL = "https://service.berlin.de/terminvereinbarung/termin/day/1714428000/"

func main() {
	for {
		poll()
		time.Sleep(time.Duration(60+rand.Intn(60)) * time.Second)
	}
}

func poll() {
	jar, _ := cookiejar.New(nil)
	httpClient := &http.Client{Jar: jar}

	// make first request to get session-cookies
	req := createGetRequest(entryURL)
	res, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// make second request to get HTML page with appointments
	req = createGetRequest(appointmentURL)
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
		req := createTelegramSendMessageRequest("Hey, I've got some free appointments. Go get em'! https://service.berlin.de/terminvereinbarung/termin/all/120686/")
		resp, err := httpClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
		}
	}
	fmt.Printf("Found %v days with open slots (%s)\n", len(bookableDataPoints.Nodes), time.Now().Format("02.01.2006 15:04:05"))
}

func createGetRequest(url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	return req
}
