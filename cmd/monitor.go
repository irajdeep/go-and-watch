package main

import (
	"log"
)

// Receiver channel from parser
// To be called from main.go
func processAndMonitor(statsData <-chan AggregatedStats) {
	// get the desired data from the channel
	monitorStat, _ := <-statsData
	monitorEndpoint(monitorStat.EndPointStats)
	monitorStatusCode(monitorStat.RequestStatusStats)

}

func monitorEndpoint(endPointStat []EndPointStat) {

	maxHits := int(0)
	totalHits := int(0)
	maxHitEndpoint := ""

	for _, element := range endPointStat {
		// index is the index where we are
		// element is the element from someSlice for where we are
		if element.Hits > maxHits {
			maxHits = element.Hits
			maxHitEndpoint = element.EndPoint
		}
		totalHits += element.Hits
	}
	log.Printf("Maximum hit endpoint %s", maxHitEndpoint)

}

func monitorStatusCode(requestStatusStats []RequestStatusStat) {
	log.Printf("Request statuscode stats over last 10 secs: %v", requestStatusStats)
}
