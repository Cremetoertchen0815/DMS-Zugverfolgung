package api

import (
	"TestApp/bl"
	"TestApp/webtool"
	"sync"
	"time"
)

const MAX_RFID_SIZE_SCANNER_DELAY_SECONDS = 50

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

	Logger chan<- webtool.LogItem
}

type SizeScannerHistoryItem struct {
	Scanner   SizeScannerId
	TimeStamp time.Time
	Train     bl.TrainId
	Value     float64
}

type RFIDScannerHistoryItem struct {
	Scanner   RFIDScannerId
	TimeStamp time.Time
	Value     bl.TrainId
}

// CreateScanners loads or creates a new Scanners object
func CreateScannerManager(logger chan<- webtool.LogItem) *ScannerManager {
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

		Logger: logger}

	return manager
}

func (manager *ScannerManager) GetRFIDFromScanner(scanner RFIDScannerId, maxTimeoutSeconds float64) (bl.TrainId, error) {
	return 187, nil
}

func (manager *ScannerManager) ContainsRFIDScanner(Id RFIDScannerId) bool {
	for _, id := range manager.RFIDScanners {
		if id == Id {
			return true
		}
	}

	return false
}
