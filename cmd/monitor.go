package main

import (
	"log"
)

// Receiver channel from parser

func processAndDisplayMonitorData(statsData <-chan AggregatedStats) {
	// get the desired data from the channel
	endPointStat, requestStatusStat := <-statsData

	maxHits := int64(0)
	totalHits := int64(0)
	maxHitEndpoint := ""

	for _, element := range endPointStat.EndPointStats {
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
