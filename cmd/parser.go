package main

import (
	"github.com/Songmu/axslogparser"
	"github.com/hpcloud/tail"
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

	go sendStatsToMonitor(monitorCh)
	go sendStatsToAlert(alertsCh)

	for lineStruct := range logCh {
		go updateDataStructure(&lineStruct)
	}
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
}

type AggregatedStats struct {
	EndPointStats []EndPointStat
	RequestStatusStats []RequestStatusStat
}

// In-memory data structures
var aggregatedStats AggregatedStats
var dataStore DataStore

func updateDataStructure(lineStruct *LogLine) {

}

func sendStatsToAlert(alertsCh chan AggregatedStats) {
	// TODO
}

func sendStatsToMonitor(monitorCh chan AggregatedStats) {
	// TODO
}

func cleanDataStore() {
	// TODO
}