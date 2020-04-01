package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gocolly/colly"
	"github.com/jordic/goics"
)

// Event is an entry in a calendar
type Event struct {
	Start, End  time.Time
	ID, Summary string
}

// Events is a collection of calendar events
type Events []*Event

const TimeLayout = "2018-02-01T13:30:01+00:00"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/feed/", feed())

	log.Print("Server started on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func feed() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/calendar")
		w.Header().Set("charset", "utf-8")
		w.Header().Set("Content-Disposition", "inline")
		w.Header().Set("filename", "calendar.ics")

		entries, err := fetchData()
		if err != nil {
			log.Fatal(err)
		}

		b := bytes.Buffer{}
		goics.NewICalEncode(&b).Encode(entries)

		writeSuccess(b.String(), w)
	})
}

func fetchData() (Events, error) {
	dividendDatesPerStock := make(map[string]string)
	result := Events{}
	PORTFOLIO := []string{"SAN.PA", "FRE.DE", "AI.PA", "BAYN.DE", "ALV.DE", "DTE.DE", "MSFT", "AIR.PA"}
	// WATCHLIST := []string{"WORK:US", "VAR1:GR", "ROG:SW", "AMZN:US", "GOOGL", "HYQ.DE", "BAS:GR", "BMW:GR", "POAHY:US", "LHA:GR", "RDSA:NA"}

	c := colly.NewCollector()

	for i := 0; i < len(PORTFOLIO); i++ {
		c.OnHTML("td[data-test]", func(e *colly.HTMLElement) {
			dataTest := e.Attr("data-test")
			if dataTest == "EX_DIVIDEND_DATE-value" {
				dividendDatesPerStock[PORTFOLIO[i]] = e.Text
			}
		})

		c.Visit("https://finance.yahoo.com/quote/" + PORTFOLIO[i])

		date, err := time.Parse(TimeLayout, dividendDatesPerStock[PORTFOLIO[i]])
		if err != nil {
			log.Fatal(err)
		}

		event := Event{
			Start:   date,
			End:     date,
			ID:      PORTFOLIO[i],
			Summary: "Some company",
		}

		result = append(result, &event)

		fmt.Println(PORTFOLIO[i] + " pays dividends on " + dividendDatesPerStock[PORTFOLIO[i]])
	}

	fmt.Println(len(PORTFOLIO), len(dividendDatesPerStock))

	return result, nil
}

// EmitICal implements the interface for goics
func (e Events) EmitICal() goics.Componenter {
	c := goics.NewComponent()
	c.SetType("VCALENDAR")
	c.AddProperty("CALSCAL", "GREGORIAN")

	for _, event := range e {
		s := goics.NewComponent()
		s.SetType("VEVENT")

		k, v := goics.FormatDateTimeField("DTSTART", event.Start)
		s.AddProperty(k, v)
		k, v = goics.FormatDateTimeField("DTEND", event.End)
		s.AddProperty(k, v)

		s.AddProperty("SUMMARY", event.Summary)
		c.AddComponent(s)
	}
	return c
}

func writeError(status int, message string, w http.ResponseWriter, err error) {
	log.Print("ERROR: ", err.Error())
	w.WriteHeader(status)
	w.Write([]byte(message))
}

func writeSuccess(message string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}
