package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	monitorCh := make(chan AggregatedStats)
	alertsCh := make(chan AggregatedStats)

	go func() {
		for stats := range monitorCh {
			jso, _ := json.Marshal(stats)
			log.Println(string(jso))
		}
	}()

	go Process(monitorCh, alertsCh)

	go processAndMonitor(monitorCh)
	go validateAndDisplayAlert(alertsCh)

	stopCh := make(chan struct{})
	defer close(stopCh)

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGTERM)
	signal.Notify(sigTerm, syscall.SIGINT)
	<-sigTerm
}
