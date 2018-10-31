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
	MaxEndPointHitThreshold: 100,
	AlertInterval:           60, // every 1 minute
}

// Receives channel from main.go
// channel should have data for alert interval , i.e 1 min
func validateAndDisplayAlert(statsData <-chan AggregatedStats) {

	alertCh := make(chan AggregatedStats)
	for aggregatedStats := range alertCh {
		go displayAlert(aggregatedStats)

		time.Sleep(tempAlertConfig.AlertInterval * time.Second)
		go computeAggregateStats(tempAlertConfig.AlertInterval, alertCh)
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
	fmt.Println("Violating endpoints:")
	fmt.Println(violatingEndpoints)
}
