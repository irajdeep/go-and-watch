package main

// Receiver channel from parser
func processMonitorData(statsData <-chan AggregatedStats) {
	// get the desired data from the channel
	endPointStat, requestStatusStat := <-statsData

	tatalHits := 0
	maxHHitEndpoint := ""

	for k, v := range endPointStat.EndPointStats {

	}

	for k, v := range endPointStat.RequestStatusStats {

	}

}
