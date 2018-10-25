package main

import "log"

// Receives channel from main.go
// channel should have data for alert interval , i.e 1 min
func validateAndDisplayAlert(statsData <-chan AggregatedStats) {
	// Should be taken from command line parameter in main.go
	alertThreshold := int(20)
	alertCount := 0

	monitorStat, _ := <-statsData

	for _, element := range monitorStat.EndPointStats {
		if element.Hits > alertThreshold {
			alertCount += 1
		}
	}
	log.Printf("Number of endpoints with unsual high traffic %d", alertCount)
}
