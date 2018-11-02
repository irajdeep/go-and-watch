package main

import (
	"time"

	"github.com/prometheus/common/log"
)

const retainDataStoreSeconds int = 130

func cleanDataStore() {
	ticker := time.NewTicker(30 * time.Second)
	quit := make(chan struct{})

	// clean entries from the datastore whose time difference is greater than 2 minutes (130 seconds some extra buffer)
	// which is maximum duration data we need right now for alert
	go func() {
		for {
			select {
			case <-ticker.C:
				// lock datastore before cleaning
				dataStore.mutex.Lock()
				log.Infof("Starting to clean datastore..")

				// current epoch since expired
				currentLengthTimeStampsSorted := len(dataStore.TimeStampsSorted)
				if currentLengthTimeStampsSorted > retainDataStoreSeconds {
					timestampsToRemove := dataStore.TimeStampsSorted[:currentLengthTimeStampsSorted-retainDataStoreSeconds]
					dataStore.TimeStampsSorted = dataStore.TimeStampsSorted[currentLengthTimeStampsSorted-retainDataStoreSeconds:]
					for _, timestamp := range timestampsToRemove {
						delete(dataStore.RequestStatusStats, timestamp)
						delete(dataStore.EndPointStats, timestamp)
						delete(dataStore.TimeStampsDict, timestamp)
					}
				}

				dataStore.mutex.Unlock()
				log.Info("...Cleaning datastore done")
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
