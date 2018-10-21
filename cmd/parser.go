package parser

import (
	"os"
	"bufio"
)

type LogLine struct {
	FormattedLine string // as is log line
	ClientIP string
	When string
	ReqType string
	EndPoint string
	StatusCode int
	Info string
}

func (logLine *LogLine) parseLine() {
	// TODO parse logic using logLine.FormattedLine

}

func ParseLogFile(fileName string) ([]LogLine) {
	linesCh := make(chan string)

	go readLogFile(fileName, linesCh)

	log := make([]LogLine, 100)
	for line := range linesCh {
		lineStruct := &LogLine{FormattedLine: line}
		lineStruct.parseLine()
		log = append(log, lineStruct)
	}

	return log
}

func readLogFile(fileName string, linesCh chan string) {
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