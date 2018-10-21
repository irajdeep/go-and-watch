package main

import (
	"os"
	"bufio"
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
	t, err := tail.TailFile(filePath, tail.Config{
		Follow: true,
		ReOpen: true})

	log := make([]LogLine, 0, 100)
	for line := range t.lines {
		lineStruct := LogLine{FormattedLine: line}
		lineStruct.parseLine()
		logCh <- lineStruct
	}
}