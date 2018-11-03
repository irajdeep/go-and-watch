package main

import (
	"fmt"
	"time"
)

type AlertConfig struct {
	// fire alert if:
	MaxEndPointHitThreshold int // a endpoint has these many or more hits
	AlertInterval           time.Duration
}

var tempAlertConfig AlertConfig = AlertConfig{
	MaxEndPointHitThreshold: 5,
	AlertInterval:           60, // every 1 minute
}

// Receives channel from main.go
// channel should have data for alert interval , i.e 1 min
func validateAndDisplayAlert() {

	alertCh := make(chan AggregatedStats)
	go func() {
		monitorTicker := time.NewTicker(tempAlertConfig.AlertInterval * time.Second)
		for {
			select {
			case <-monitorTicker.C:
				go computeAggregateStats(tempAlertConfig.AlertInterval, alertCh)
			default:

			}
		}
	}()

	for aggregatedStats := range alertCh {
		displayAlert(aggregatedStats)
	}
}

func displayAlert(aggregatedStats AggregatedStats) {

	violatingEndpoints := make([]EndPointStat, 0, 100)
	for _, endpointHits := range aggregatedStats.EndPointStats {
		if endpointHits.Hits > tempAlertConfig.MaxEndPointHitThreshold {
			violatingEndpoints = append(violatingEndpoints, endpointHits)
		}
	}

	fmt.Println("****ALERT****")
	fmt.Println("Violating endpoints:", violatingEndpoints)
}
