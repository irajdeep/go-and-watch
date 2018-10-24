package main

import (
	"log"
	"sync"

	"github.com/Songmu/axslogparser"
	"github.com/hpcloud/tail"
)

type LogLine struct {
	FormattedLine string // as is log line
	parsedLog     *axslogparser.Log
}

func (logLine *LogLine) parseLine() {
	parsedLog, err := axslogparser.Parse(logLine.FormattedLine)
	log.Println(parsedLog, err)
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

func Process(monitorCh chan AggregatedStats, alertsCh chan AggregatedStats) {
	logCh := make(chan LogLine)
	go ParseLogFile("../sample-log/sample.log", logCh)

	aggregatedStatsCh := make(chan AggregatedStats)

	for lineStruct := range logCh {
		go updateDataStructure(lineStruct, aggregatedStatsCh)
		go computeAggregatedStatsAndSend(aggregatedStatsCh)
	}

	go sendStatsToMonitor(monitorCh, aggregatedStatsCh)
	go sendStatsToAlerts(alertsCh, aggregatedStatsCh)

	go cleanDataStore()
}

type EndPointStat struct {
	EndPoint string
	hits     int
}

type RequestStatusStat struct {
	Status int
	count  int
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
var aggregatedStats AggregatedStats
var dataStore DataStore

func updateDataStructure(lineStruct LogLine, aggregatedStatsCh chan AggregatedStats) {
	parsedLog := lineStruct.parsedLog
	epoch := parsedLog.Time.Unix()

	_, exists := dataStore.TimeStampsDict[epoch]

	// make space for new timestamp
	if !exists {
		dataStore.mutex.Lock()

		dataStore.TimeStampsDict[epoch] = true
		dataStore.TimeStampsSorted = append(dataStore.TimeStampsSorted, epoch)
		dataStore.EndPointStats[epoch] = make(map[string]int)
		dataStore.RequestStatusStats[epoch] = make(map[int]int)

		dataStore.mutex.Unlock()
	}

	go actualUpdateDataStructure(lineStruct, aggregatedStatsCh)
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

func computeAggregatedStatsAndSend(aggregatedStatsCh chan AggregatedStats) {
	dataStore.mutex.Lock()
	defer dataStore.mutex.Unlock()

	timeStamps := dataStore.TimeStampsSorted
	leftIdx := len(timeStamps) - 10
	if leftIdx <= 0 {
		leftIdx = 0
	}

	// pick last 10 epochs(seconds)
	windowStart := len(timeStamps) - 10
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

	var aggregatedStats AggregatedStats
	for uri, count := range seenURIs {
		aggregatedStats.EndPointStats = append(aggregatedStats.EndPointStats, EndPointStat{EndPoint: uri, hits: count})
	}
	for status, count := range seenStatus {
		aggregatedStats.RequestStatusStats = append(aggregatedStats.RequestStatusStats, RequestStatusStat{Status: status, count: count})
	}

	aggregatedStatsCh <- aggregatedStats
}

func sendStatsToMonitor(monitorCh chan AggregatedStats, aggregatedStats chan AggregatedStats) {
	// TODO add non blocking select
	for stats := range aggregatedStats {
		monitorCh <- stats
	}
}

func sendStatsToAlerts(alertsCh chan AggregatedStats, aggregatedStats chan AggregatedStats) {
	// TODO add non blocking select
	for stats := range aggregatedStats {
		alertsCh <- stats
	}
}

func cleanDataStore() {
	// TODO
}
