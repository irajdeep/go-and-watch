package main

import (
	"testing"
)

func TestParseLogFile(t * testing.T) {
	fileName := "../sample-log/sample.log"
	log := ParseLogFile(fileName)
	for _, l := range log {
		t.Log(l.parsedLog)
	}
	t.Log(len(log))
	t.Log("asf")
}