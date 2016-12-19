package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

// MeetupEvent is a basic event as provided by the Meetup API.
type MeetupEvent struct {
	Created       int64  `json:"created"`
	ID            string `json:"id"`
	Name          string `json:"name"`
	Status        string `json:"status"`
	Time          int64  `json:"time"`
	Updated       int64  `json:"updated"`
	UtcOffset     int    `json:"utc_offset"`
	WaitlistCount int    `json:"waitlist_count"`
	YesRsvpCount  int    `json:"yes_rsvp_count"`
	Venue         struct {
		ID                   int     `json:"id"`
		Name                 string  `json:"name"`
		Lat                  float64 `json:"lat"`
		Lon                  float64 `json:"lon"`
		Repinned             bool    `json:"repinned"`
		Address1             string  `json:"address_1"`
		City                 string  `json:"city"`
		Country              string  `json:"country"`
		LocalizedCountryName string  `json:"localized_country_name"`
	} `json:"venue"`
	Group struct {
		Created  int64   `json:"created"`
		Name     string  `json:"name"`
		ID       int     `json:"id"`
		JoinMode string  `json:"join_mode"`
		Lat      float64 `json:"lat"`
		Lon      float64 `json:"lon"`
		Urlname  string  `json:"urlname"`
		Who      string  `json:"who"`
	} `json:"group"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Visibility  string `json:"visibility"`
}

func queryMeetupEvents(group string, status string) []MeetupEvent {
	url := "https://api.meetup.com/SciFiMuc/events?status=past"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Could not create request:", err)
		return nil
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Could not perform request:", err)
		return nil
	}
	defer resp.Body.Close()

	var events []MeetupEvent

	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		log.Println(err)
	}
	return events
}

func writeAsCsv(events []MeetupEvent) {

	file, err := os.Create("test.csv")
	if err != nil {
		log.Fatalln("error creating file", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Write([]string{"name", "time"})

	for _, event := range events {
		var secondsSinceEpoch int64
		secondsSinceEpoch = event.Time / 1000
		t := time.Unix(secondsSinceEpoch, 0)
		writer.Write([]string{event.Name, t.Format("2006-01-02 15:04")})
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
}

func main() {
	log.Println("Get data from Meetup API and Write to CSV")
	events := queryMeetupEvents("", "")
	writeAsCsv(events)
	log.Println("Done.")
}
