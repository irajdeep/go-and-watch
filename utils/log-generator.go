package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/brianvoe/gofakeit"
)

const (
	ApacheCommonLog = "%s - - [%s] \"%s %s\" %d %d"
)

var endPoints = [...]string{"foo",
	"bar",
	"/",
	"/admin",
	"/quora/?id=1",
	"/get_room_page",
	"get_foo_bar",
	"getIdentifier",
	"get_meaning_of_42"}

// sample log line
// 77.179.66.156 - - [25/Oct/2016:14:49:33 +0200] "GET / HTTP/1.1" 200 612 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.59 Safari/537.36"
func main() {

	filePath := flag.String("log-file-path", "/tmp/access.log", "Pass the filepath to write logs to")
	flag.Parse()

	var f *os.File

	if _, err := os.Stat(*filePath); os.IsNotExist(err) {
		// path/to/whatever does not exist
		f, err = os.Create(*filePath)
		if err != nil {
			log.Printf("Failed to create file %v", err)
		}
	} else {
		f, err = os.OpenFile(*filePath, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
	}

	defer f.Close()

	// Default number of times the log f gets written per second
	timesPerSecond := 10
	for {
		logWriteData := newApacheCommonLog()
		log.Printf("Writting .... %s", logWriteData)
		_, err := f.WriteString(logWriteData + "\n")

		if err != nil {
			log.Fatalf("Failed to write to f : %v", err)
		}

		time.Sleep(time.Duration(1e9 / timesPerSecond)) //
	}
}

// "%s - - [%s] \"%s %s\" %d %d"
func newApacheCommonLog() string {
	return fmt.Sprintf(
		ApacheCommonLog,
		gofakeit.IPv4Address(),
		fakeFormattedCurrentTime(),
		gofakeit.HTTPMethod(),
		randResourceURI(),
		gofakeit.StatusCode(),
		gofakeit.Number(0, 30000),
	)
}

// 11/Jun/2017:05:56:04 +0900
func fakeFormattedCurrentTime() string {
	t := time.Now()
	return fmt.Sprintf("%02d/%s/%02d:%02d:%02d:%02d +0200",
		t.Day(), t.Month().String()[:3], t.Year(),
		t.Hour(), t.Minute(), t.Second())
}

func randResourceURI() string {
	num := gofakeit.Number(0, len(endPoints)-1)
	return endPoints[num] + " HTTP/1.1"
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
func updateTimesPerSecond() {

}
