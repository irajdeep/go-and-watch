package main

import (
	"fmt"
	"time"
)

type EndPointStat struct {
	EndPoint string
	Hits     int
}

func (endpointStat *EndPointStat) Stringer() string {
	return fmt.Sprintf("%s hits: %d", endpointStat.EndPoint, endpointStat.Hits)
}

type RequestStatusStat struct {
	Status int
	Count  int
}

type AggregatedStats struct {
	EndPointStats      []EndPointStat
	RequestStatusStats []RequestStatusStat
}

func computeAggregateStats(duration time.Duration, aggregatedStatsCh chan AggregatedStats) {
	dataStore.mutex.Lock()
	defer dataStore.mutex.Unlock()

	timeStamps := dataStore.TimeStampsSorted

	// pick last duration number of seconds
	windowStart := len(timeStamps) - int(duration)
	if windowStart < 0 {
		windowStart = 0
	}
	timeStamps = timeStamps[windowStart:]

	seenURIs := make(map[string]int)
	seenStatus := make(map[int]int)
	for _, t := range timeStamps {
		for uri, count := range dataStore.EndPointStats[t] {
			seenURIs[uri] += count
		}
		for status, count := range dataStore.RequestStatusStats[t] {
			seenStatus[status] += count
		}
	}

	aggregatedStats := AggregatedStats{
		EndPointStats:      make([]EndPointStat, 0, 100),
		RequestStatusStats: make([]RequestStatusStat, 0, 100),
	}
	for uri, count := range seenURIs {
		aggregatedStats.EndPointStats = append(aggregatedStats.EndPointStats, EndPointStat{EndPoint: uri, Hits: count})
	}
	for status, count := range seenStatus {
		aggregatedStats.RequestStatusStats = append(aggregatedStats.RequestStatusStats, RequestStatusStat{Status: status, Count: count})
	}

	aggregatedStatsCh <- aggregatedStats
}
