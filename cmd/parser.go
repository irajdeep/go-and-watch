package parser

import (
	"os"
	"bufio"
	"github.com/Songmu/axslogparser"
)

type LogLine struct {
	FormattedLine string // as is log line
	parsedLog *axslogparser.Log
}

func (logLine *LogLine) parseLine() {
	parsedLog, _ := axslogparser.Parse(logLine.FormattedLine)
	logLine.parsedLog = parsedLog
}

func ParseLogFile(filePath string) ([]LogLine) {
	linesCh := make(chan string)

	go readLogFile(filePath, linesCh)

	log := make([]LogLine, 0, 100)
	for line := range linesCh {
		lineStruct := LogLine{FormattedLine: line}
		lineStruct.parseLine()
		log = append(log, lineStruct)
	}

	return log
}

func readLogFile(filePath string, linesCh chan string) {
	defer close(linesCh)

	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		linesCh <- scanner.Text()
	}
	err = scanner.Err()
}