package main

import (
	"bytes"
	"log"
	"net/http"
	"strings"
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
type Events []Event

// TimeLayout is what we want the iCal date format to look like
const TimeLayout = "Jan 02, 2006"

// PORTFOLIO for developing
var PORTFOLIO = []string{"SAN.PA", "FRE.DE", "AI.PA", "BAYN.DE", "ALV.DE", "DTE.DE", "MSFT", "AIR.PA"}

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

		queryParams := r.URL.Query()
		symbols := strings.Split(queryParams.Get("symbols"), ",")

		entries, err := FetchData(symbols)
		if err != nil {
			log.Fatal(err)
		}

		b := bytes.Buffer{}
		goics.NewICalEncode(&b).Encode(entries)

		writeSuccess(b.String(), w)
	})
}

// func search() http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		symbols, err := FindSymbol()

// 		writeSuccess()
// 	})
// }

// FindSymbol gets a list of symbols from the yahoo lookup page
// func FindSymbol(search string) ([]string, error) {
// 	var symbols []string
// 	c := colly.NewCollector()
// 	c.OnHTML("h1", func(e *colly.HTMLElement) {
// 		event.Summary = "Ex Dividend Date: " + e.Text
// 	})
// 	c.Visit("https://finance.yahoo.com/lookup/equity?s=" + search)

// 	return symbols, nil
// }

// FetchData gets the ex dividend dates from the yahoo finance page
func FetchData(portfolio []string) (Events, error) {
	result := Events{}
	c := colly.NewCollector()
	for i := 0; i < len(portfolio); i++ {
		var event Event
		c.OnHTML("td[data-test]", func(e *colly.HTMLElement) {
			dataTest := e.Attr("data-test")
			if dataTest == "EX_DIVIDEND_DATE-value" {
				date, err := time.Parse(TimeLayout, e.Text)
				if err != nil {
					log.Fatal(err)
				}
				event.ID = portfolio[i]
				event.Start = date
				event.End = date
			}
		})
		c.OnHTML("h1", func(e *colly.HTMLElement) {
			event.Summary = "Ex Dividend Date: " + e.Text
		})
		c.Visit("https://finance.yahoo.com/quote/" + portfolio[i])
		result = append(result, event)
	}
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
