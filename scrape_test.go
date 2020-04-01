package main

import (
	"fmt"
	"testing"
)

// WATCHLIST := []string{"WORK:US", "VAR1:GR", "ROG:SW", "AMZN:US", "GOOGL", "HYQ.DE", "BAS:GR", "BMW:GR", "POAHY:US", "LHA:GR", "RDSA:NA"}
// PORTFOLIO := []string{"SAN.PA", "FRE.DE", "AI.PA", "BAYN.DE", "ALV.DE", "DTE.DE", "MSFT", "AIR.PA"}

func TestFetchData(t *testing.T) {
	var mock = struct {
		portfolio []string
		dates     map[string]string
		summaries map[string]string
	}{
		[]string{"SAN.PA", "FRE.DE", "ALV.DE"},
		map[string]string{
			"SAN.PA": "2020-05-04 00:00:00 +0000 UTC",
			"FRE.DE": "2020-05-21 00:00:00 +0000 UTC",
			"ALV.DE": "2020-05-07 00:00:00 +0000 UTC",
		},
		map[string]string{
			"SAN.PA": "Ex Dividend Date: SAN.PA - Sanofi",
			"FRE.DE": "Ex Dividend Date: FRE.DE - Fresenius SE & Co. KGaA",
			"ALV.DE": "Ex Dividend Date: ALV.DE - Allianz SE",
		},
	}

	t.Run("Get company name and ex dividend date", func(t *testing.T) {
		events, err := FetchData(mock.portfolio)
		if err != nil {
			t.Error((err))
		}

		for _, event := range events {
			fmt.Println(event)
			if mock.dates[event.ID] != event.Start.String() {
				t.Errorf("got %v, want %v", event.Start.String(), mock.dates[event.ID])
			}
			if mock.summaries[event.ID] != event.Summary {
				t.Errorf("got %v, want %v", event.Summary, mock.summaries[event.ID])
			}
		}
	})
}
