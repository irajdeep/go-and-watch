package main

import (
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
			log.Println(stats)
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
