package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	monitorCh := make(chan AggregatedStats)
	alertsCh := make(chan AggregatedStats)

	go ProcessLogs(monitorCh, alertsCh)

	go processAndMonitor(monitorCh)
	go validateAndDisplayAlert(alertsCh)

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGTERM)
	signal.Notify(sigTerm, syscall.SIGINT)
	<-sigTerm
}
