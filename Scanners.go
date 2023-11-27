package main

import (
	"sync"
	"time"
)

type SizeScannerId int
type RFIDScannerId int

type ScannerManager struct {
	SizeScanners       map[SizeScannerId]RFIDScannerId
	SizeScannerHistory map[int]SizeScannerHistoryItem
	SizeScanned        chan SizeScannerHistoryItem
	SizeScannerMutex   sync.Mutex

	RFIDScanners       []RFIDScannerId
	RFIDScannerHistory map[int]RFIDScannerHistoryItem
	RFIDScanned        chan RFIDScannerHistoryItem
	RFIDScannerMutex   sync.Mutex
}

type SizeScannerHistoryItem struct {
	Scanner   SizeScannerId
	TimeStamp time.Time
	Train     TrainId
	Value     float64
}

type RFIDScannerHistoryItem struct {
	Scanner   RFIDScannerId
	TimeStamp time.Time
	Value     TrainId
}

// CreateScanners loads or creates a new Scanners object
func CreateScannerManager() *ScannerManager {
	sizeChannel := make(chan SizeScannerHistoryItem)
	rfidChannel := make(chan RFIDScannerHistoryItem)

	//Create a new scanner manager(TODO: Load from file)
	manager := &ScannerManager{
		SizeScanners:       make(map[SizeScannerId]RFIDScannerId),
		SizeScannerHistory: make(map[int]SizeScannerHistoryItem),
		SizeScanned:        sizeChannel,
		SizeScannerMutex:   sync.Mutex{},

		RFIDScanners:       make([]RFIDScannerId, 0),
		RFIDScannerHistory: make(map[int]RFIDScannerHistoryItem),
		RFIDScanned:        rfidChannel,
		RFIDScannerMutex:   sync.Mutex{},
	}

	return manager
}
