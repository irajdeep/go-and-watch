package main

import (
	"log"
)

// Receiver channel from parser

func processAndMonitor(statsData <-chan AggregatedStats) {
	// get the desired data from the channel
	endPointStat, requestStatusStat := <-statsData
	monitorEndpoint(endPointStat.EndPointStats)
	monitorStatusCode(endPointStat.RequestStatusStats)

}

func monitorEndpoint(endPointStat []EndPointStat) {

	maxHits := int64(0)
	totalHits := int64(0)
	maxHitEndpoint := ""

	for _, element := range endPointStat {
		// index is the index where we are
		// element is the element from someSlice for where we are
		if element.hits > maxHits {
			maxHits = element.hits
			maxHitEndpoint = element.EndPoint
		}
		totalHits += element.hits
	}
	log.Printf("Maximum hit endpoint %s", maxHitEndpoint)

}

func monitorStatusCode(requestStatusStats []RequestStatusStat) {

}
