package main

func main() {
	monitorCh := make(chan AggregatedStats)
	alertsCh := make(chan AggregatedStats)

	go Process(monitorCh, alertsCh)

	go processAndMonitor(monitorCh)
	go validateAndDisplayAlert(alertsCh)
}
