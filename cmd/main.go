package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {

	go ProcessLogs()

	go processAndMonitor()
	go validateAndDisplayAlert()

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGTERM)
	signal.Notify(sigTerm, syscall.SIGINT)
	<-sigTerm
}
