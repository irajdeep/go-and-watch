package main

import (
	"github.com/Songmu/axslogparser"
	"github.com/hpcloud/tail"
	"sync"
)

type LogLine struct {
	FormattedLine string // as is log line
	parsedLog *axslogparser.Log
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

func Process(monitorCh chan AggregatedStats, alertsCh chan AggregatedStats) {
	logCh := make(chan LogLine)
	go ParseLogFile("../sample-log/sample.log", logCh)

	aggregatedStatsCh := make(chan AggregatedStats)
	
	for lineStruct := range logCh {
		go updateDataStructure(&lineStruct, aggregatedStatsCh)
		go computeAggregatedStatsAndSend(aggregatedStatsCh)
	}

	go sendStatsToMonitor(monitorCh, aggregatedStatsCh)
	go sendStatsToAlerts(alertsCh, aggregatedStatsCh)

	go cleanDataStore()
}

type EndPointStat struct {
	EndPoint string
	hits int64
}

type RequestStatusStat struct {
	Status int
	count int64
}

type DataStore struct {
	EndPointStats map[string]map[string]int // timeStamp -> endpoint->count
	RequestStatusStats map[string]map[int]int
	TimeStampsSorted []string
	TimeStampsDict map[string]bool

	mutex sync.Mutex
}

type AggregatedStats struct {
	EndPointStats []EndPointStat
	RequestStatusStats []RequestStatusStat
}

// In-memory data structures
var aggregatedStats AggregatedStats
var dataStore DataStore

func updateDataStructure(lineStruct *LogLine, aggregatedStatsCh chan AggregatedStats) {
	parsedLog := lineStruct.parsedLog
	_, exists := dataStore.TimeStampsDict[parsedLog.TimeStr]

	// make space for new timestamp
	if !exists {
		dataStore.mutex.Lock()

		dataStore.TimeStampsDict[parsedLog.TimeStr] = true
		dataStore.TimeStampsSorted = append(dataStore.TimeStampsSorted, parsedLog.TimeStr)
		dataStore.EndPointStats[parsedLog.TimeStr] = make(map[string]int)
		dataStore.RequestStatusStats[parsedLog.TimeStr] = make(map[int]int)

		dataStore.mutex.Unlock()
	}

	go actualUpdateDataStructure(lineStruct, aggregatedStatsCh)
}

func actualUpdateDataStructure(lineStruct *LogLine, aggregatedStatsCh chan AggregatedStats) {
	parsedLog := lineStruct.parsedLog
	_, exists := dataStore.TimeStampsDict[parsedLog.TimeStr]
	if !exists {
		return
	}

	dataStore.mutex.Lock()
	
	endPointStats := dataStore.EndPointStats[parsedLog.TimeStr]
	requestStatusStats := dataStore.RequestStatusStats[parsedLog.TimeStr]
	endPointStats[parsedLog.RequestURI]++
	requestStatusStats[parsedLog.Status]++

	dataStore.mutex.Unlock()
}

func computeAggregatedStatsAndSend(aggregatedStatsCh chan AggregatedStats) {
	dataStore.mutex.Lock()
	defer dataStore.mutex.Unlock()


	timeStamps := dataStore.TimeStampsSorted
	timeStamps = timeStamps[len(timeStamps) - 10 : ] // pick last 10 time stamps. FIXME: for last 10 seconds
	
	var aggregatedStats AggregatedStats
	

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