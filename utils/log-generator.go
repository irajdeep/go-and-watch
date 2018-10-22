package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"
)

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

	ips := []string{"77.179.66.156",
		"127.0.0.1",
		"127.0.0.2",
		"134.23.22.11",
		"110.21.22.1"}

	requestType := []string{"GET",
		"POST",
		"PUT",
	}

	endPoints := []string{"foo",
		"bar",
		"/",
		"/admin",
		"/quora/?id=1",
		"/get_room_page",
		"get_foo_bar",
		"getIdentifier",
		"get_meaning_of_42"}

	statusCode := []string{"200", "201", "404", "400", "500"}

	// Constant log parts
	//uselessRequestInfoLine := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.59 Safari/537.36"
	protocolType := "HTTP/1.1"

	// Default number of times the log f gets written per second
	timesPerSecond := 10
	for {
		// Write to the f
		ipIndex := randInt(0, len(ips)-1)
		requestTypeIndex := randInt(0, len(requestType)-1)
		endPointsIndex := randInt(0, len(endPoints)-1)
		statusCodeIndex := randInt(0, len(statusCode)-1)
		currentTime := time.Now().Format("25/Oct/2016:14:49:33 +0200")

		logWriteData := ips[ipIndex] + " - - " + "[" + currentTime + "]" +
			" \"" + requestType[requestTypeIndex] + " " + endPoints[endPointsIndex] + " " + protocolType +
			"\" " + statusCode[statusCodeIndex] + " 301" + " \"-\"" +
			" \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.59 Safari/537.36\""

		log.Printf("Writting .... %s", logWriteData)
		_, err := f.WriteString(logWriteData + "\n")

		if err != nil {
			log.Fatalf("Failed to write to f : %v", err)

		}

		time.Sleep(time.Duration(1e9 / timesPerSecond)) //
	}
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
func updateTimesPerSecond() {

}
