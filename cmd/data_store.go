package main

import "sync"

type DataStore struct {
	EndPointStats      map[int64]map[string]int // timeStamp -> endpoint->count
	RequestStatusStats map[int64]map[int]int
	TimeStampsSorted   []int64 // number of seconds
	TimeStampsDict     map[int64]bool

	mutex sync.Mutex
}

// In-memory data structure
var dataStore DataStore

func initDataStore() {
	dataStore.EndPointStats = make(map[int64]map[string]int)
	dataStore.RequestStatusStats = make(map[int64]map[int]int)
	dataStore.TimeStampsSorted = make([]int64, 0, 100)
	dataStore.TimeStampsDict = make(map[int64]bool)

	// Clean data store periodically
	go cleanDataStore()
}

func updateDataStructure(lineStruct LogLine) {
	parsedLog := lineStruct.parsedLog
	epoch := parsedLog.Time.Unix()

	dataStore.mutex.Lock()
	_, exists := dataStore.TimeStampsDict[epoch]
	// make space for new timestamp
	if !exists {
		dataStore.TimeStampsDict[epoch] = true
		dataStore.TimeStampsSorted = append(dataStore.TimeStampsSorted, epoch)
		dataStore.EndPointStats[epoch] = make(map[string]int)
		dataStore.RequestStatusStats[epoch] = make(map[int]int)
	}
	dataStore.mutex.Unlock()

	actualUpdateDataStructure(lineStruct)
}

func actualUpdateDataStructure(lineStruct LogLine) {
	parsedLog := lineStruct.parsedLog
	epoch := parsedLog.Time.Unix()
	_, exists := dataStore.TimeStampsDict[epoch]
	if !exists {
		return
	}

	dataStore.mutex.Lock()

	endPointStats := dataStore.EndPointStats[epoch]
	requestStatusStats := dataStore.RequestStatusStats[epoch]
	endPointStats[parsedLog.RequestURI]++
	requestStatusStats[parsedLog.Status]++

	dataStore.mutex.Unlock()
}
