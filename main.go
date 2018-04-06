package main

import (
	"net/http"
	"bytes"
	"os"
	"log"
	"github.com/PuerkitoBio/goquery"
)

import _ "github.com/joho/godotenv/autoload"

func getSoups() []string {
	res, err := http.Get("http://pksoup.com/locations/covent-garden-market")
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

	var soups []string
	doc.Find("#todays-soups-inner .soup-menu-item").Each(func(i int, s *goquery.Selection) {
		soup := s.Find("h2").Text()
		soups = append(soups, soup)
	})

	return soups
}

func buildPayload(soups []string) bytes.Buffer {
	var payload bytes.Buffer
	payload.WriteString("{\"text\":\"")
	for _, soup := range soups {
		payload.WriteString(soup)
		payload.WriteString("\\n")
	}
	payload.WriteString("\"}")

	return payload;
}

func sendRequest(payload bytes.Buffer) {
	endpoint := os.Getenv("SLACK_ENDPOINT")

	req, err := http.NewRequest("POST", endpoint, &payload)
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func main() {
	payload := buildPayload(getSoups())
	sendRequest(payload)
}
