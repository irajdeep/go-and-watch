package main

import (
	"log"
	"time"
)

type MonitorSetting struct {
	Interval int64 // in seconds
}

// Configurable
var monitorSettings MonitorSetting = MonitorSetting{Interval: int64(3)}

// Receiver channel from parser
// To be called from main.go
func processAndMonitor(statsData <-chan AggregatedStats) {
	startTime := time.Now().Unix()

	lastMoniteredTime := startTime
	for {
		select {
		case monitorStat := <-statsData:
			currentTime := time.Now().Unix()
			if currentTime >= lastMoniteredTime+monitorSettings.Interval {
				go monitorEndpoint(monitorStat.EndPointStats)
				go monitorStatusCode(monitorStat.RequestStatusStats)
				lastMoniteredTime = currentTime

				time.Sleep(1000 * time.Millisecond)
			}
		default:
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

func monitorEndpoint(endPointStat []EndPointStat) {
	maxHits := int(0)
	totalHits := int(0)
	maxHitEndpoint := ""

	for _, element := range endPointStat {
		if element.Hits > maxHits {
			maxHits = element.Hits
			maxHitEndpoint = element.EndPoint
		}
		totalHits += element.Hits
	}
	log.Printf("Total hits %d", totalHits)
	log.Printf("Maximum hit endpoint %s hits: %d", maxHitEndpoint, maxHits)
}

func monitorStatusCode(requestStatusStats []RequestStatusStat) {
	log.Printf("Request statuscode stats over last %d secs: %v", monitorSettings.Interval, requestStatusStats)
}
