package main

import (
	"github.com/Songmu/axslogparser"
	"github.com/hpcloud/tail"
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

	go func() {
		for lineStruct := range logCh {
			updateDataStructure(lineStruct)
		}
	}()

}
