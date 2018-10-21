package parser

import (
	"testing"
	"fmt"
)

func TestParseLogFile(t * testing.T) {
	fileName := "../sample-log/sample.log"
	log := ParseLogFile(fileName)
	for l := range log {
		fmt.Println("%s", l.FormattedLine)
	}
}