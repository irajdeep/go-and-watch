package main

func main() {
	monitorCh := make(chan AggregatedStats)
	alertsCh := make(chan AggregatedStats)

	go Process(monitorCh, alertsCh)

	go processAndMonitor(monitorCh)
	go validateAndDisplayAlert(alertsCh)

	/**

	stopCh := make(chan struct{})
	defer close(stopCh)

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGTERM)
	signal.Notify(sigTerm, syscall.SIGINT)
	<-sigTerm
	 */
}
