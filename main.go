package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

var entryURL = "https://service.berlin.de/terminvereinbarung/termin/all/120686/"
var appointmentURL = "https://service.berlin.de/terminvereinbarung/termin/day/1714428000/"

func main() {
	cookies := make([]*http.Cookie, 0)
	/*
		cookies = append(cookies, &http.Cookie{Name: "zmsappointment-session", Value: "inProgress"})
		cookies = append(cookies, &http.Cookie{Name: "Zmsappointment", Value: "nhbr3g3lmklo925mvto4m8pg7j"})
		cookies = append(cookies, &http.Cookie{Name: "wt_rla", Value: "102571513503709%2C64%2C1712675202268"})
	*/
	baseURL, _ := url.Parse("https://service.berlin.de")
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(baseURL, cookies)
	httpClient := &http.Client{
		Jar: jar,
	}
	httpClient.Jar.SetCookies(baseURL, cookies)

	req, err := http.NewRequest("GET", entryURL, nil)
	if err != nil {
		log.Fatalln(err)
	}
	addHeaders(req)

	res, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	req, err = http.NewRequest("GET", appointmentURL, nil)
	if err != nil {
		log.Fatalln(err)
	}
	addHeaders(req)
	res, err = httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

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
		fmt.Print("Sucess! ")
	}
	fmt.Printf("Found %v days with open slots\n", len(bookableDataPoints.Nodes))
}

func addHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
}
