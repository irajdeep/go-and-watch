package main

import (
	"sync"
	"time"

	"github.com/Songmu/axslogparser"
	"github.com/hpcloud/tail"
	"github.com/prometheus/common/log"
)

type LogLine struct {
	FormattedLine string // as is log line
	parsedLog     *axslogparser.Log
}

func (logLine *LogLine) parseLine() {
	parsedLog, _ := axslogparser.Parse(logLine.FormattedLine)
	logLine.parsedLog = parsedLog
}

func ParseLogFile(filePath string, logCh chan LogLine) {
	t, _ := tail.TailFile(filePath, tail.Config{
		Follow: true,
		ReOpen: true})

	for line := range t.Lines {
		lineStruct := LogLine{FormattedLine: line.Text}
		lineStruct.parseLine()
		logCh <- lineStruct
	}
}

func ProcessLogs() {
	logCh := make(chan LogLine)
	go ParseLogFile("../sample-log/sample.log", logCh)

	initDataStore()

	aggregatedStatsCh := make(chan AggregatedStats)
	go func() {
		for lineStruct := range logCh {
			updateDataStructure(lineStruct, aggregatedStatsCh)
		}
	}()

	go cleanDataStore()
}

type EndPointStat struct {
	EndPoint string
	Hits     int
}

type RequestStatusStat struct {
	Status int
	Count  int
}

type DataStore struct {
	EndPointStats      map[int64]map[string]int // timeStamp -> endpoint->count
	RequestStatusStats map[int64]map[int]int
	TimeStampsSorted   []int64 // number of seconds
	TimeStampsDict     map[int64]bool

	mutex sync.Mutex
}

type AggregatedStats struct {
	EndPointStats      []EndPointStat
	RequestStatusStats []RequestStatusStat
}

// In-memory data structures
var dataStore DataStore

func initDataStore() {
	dataStore.EndPointStats = make(map[int64]map[string]int)
	dataStore.RequestStatusStats = make(map[int64]map[int]int)
	dataStore.TimeStampsSorted = make([]int64, 0, 100)
	dataStore.TimeStampsDict = make(map[int64]bool)
}

func updateDataStructure(lineStruct LogLine, aggregatedStatsCh chan AggregatedStats) {
	parsedLog := lineStruct.parsedLog
	epoch := parsedLog.Time.Unix()

	dataStore.mutex.Lock()
	_, exists := dataStore.TimeStampsDict[epoch]
	// make space for new timestamp
	if !exists {
		dataStore.TimeStampsDict[epoch] = true
		dataStore.TimeStampsSorted = append(dataStore.TimeStampsSorted, epoch)
		dataStore.EndPointStats[epoch] = make(map[string]int)
		dataStore.RequestStatusStats[epoch] = make(map[int]int)
	}
	dataStore.mutex.Unlock()

	actualUpdateDataStructure(lineStruct, aggregatedStatsCh)
}

func actualUpdateDataStructure(lineStruct LogLine, aggregatedStatsCh chan AggregatedStats) {
	parsedLog := lineStruct.parsedLog
	epoch := parsedLog.Time.Unix()
	_, exists := dataStore.TimeStampsDict[epoch]
	if !exists {
		return
	}

	dataStore.mutex.Lock()

	endPointStats := dataStore.EndPointStats[epoch]
	requestStatusStats := dataStore.RequestStatusStats[epoch]
	endPointStats[parsedLog.RequestURI]++
	requestStatusStats[parsedLog.Status]++

	dataStore.mutex.Unlock()
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

const retainDataStoreSeconds int = 130

func cleanDataStore() {
	ticker := time.NewTicker(30 * time.Second)
	quit := make(chan struct{})

	// clean entries from the datastore whose time difference is greater than 2 minutes (130 seconds some extra buffer)
	// which is maximum duration data we need right now for alert
	go func() {
		for {
			select {
			case <-ticker.C:
				// lock datastore before cleaning
				dataStore.mutex.Lock()
				log.Infof("Starting to clean datastore..")

				// current epoch since expired
				currentLengthTimeStampsSorted := len(dataStore.TimeStampsSorted)
				if currentLengthTimeStampsSorted > retainDataStoreSeconds {
					timestampsToRemove := dataStore.TimeStampsSorted[:currentLengthTimeStampsSorted-retainDataStoreSeconds]
					dataStore.TimeStampsSorted = dataStore.TimeStampsSorted[currentLengthTimeStampsSorted-retainDataStoreSeconds:]
					for _, timestamp := range timestampsToRemove {
						delete(dataStore.RequestStatusStats, timestamp)
						delete(dataStore.EndPointStats, timestamp)
						delete(dataStore.TimeStampsDict, timestamp)
					}
				}

				dataStore.mutex.Unlock()
				log.Info("...Cleaning datastore done")
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
