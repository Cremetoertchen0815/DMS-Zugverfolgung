package main

import (
	"log"
	"net/http"
	"strconv"
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

func (manager *ScannerManager) GetRFIDFromScanner(scanner RFIDScannerId, maxTimeoutSeconds float64) (TrainId, error) {
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

// HandleSizeRequest handles all requests to the size scanner
func (manager *ScannerManager) HandleSizeRequest(w http.ResponseWriter, r *http.Request) {
	manager.SizeScannerMutex.Lock()
	defer manager.SizeScannerMutex.Unlock()

	switch r.Method {
	case "GET": //Register a new size scanner
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "No id specified", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid id specified", http.StatusBadRequest)
			return
		}

		rfidScannerStr := r.URL.Query().Get("rfidScannerId")
		if rfidScannerStr == "" {
			http.Error(w, "No rfidScannerId specified", http.StatusBadRequest)
			return
		}

		rfidScannerId, err := strconv.Atoi(rfidScannerStr)
		if err != nil || !manager.ContainsRFIDScanner(RFIDScannerId(rfidScannerId)) {
			http.Error(w, "Invalid rfidScannerId specified", http.StatusBadRequest)
			return
		}

		manager.SizeScanners[SizeScannerId(id)] = RFIDScannerId(rfidScannerId)

	case "POST": //Add a new size scan
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "No id specified", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid id specified", http.StatusBadRequest)
			return
		}

		valueStr := r.URL.Query().Get("value")
		if valueStr == "" {
			http.Error(w, "No value specified", http.StatusBadRequest)
			return
		}

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			http.Error(w, "Invalid value specified", http.StatusBadRequest)
			return
		}

		//TODO: Get the RFID scanner id from the size scanner id
		trainId, err := manager.GetRFIDFromScanner(manager.SizeScanners[SizeScannerId(id)], MAX_RFID_SIZE_SCANNER_DELAY_SECONDS)
		if err != nil {
			log.Fatalln("Error matching trains to size scanner: ", err)
			return
		}

		manager.SizeScanned <- SizeScannerHistoryItem{SizeScannerId(id), time.Now(), trainId, value}

		w.WriteHeader(http.StatusOK)
	case "DELETE": //Remove a size scanner
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "No id specified", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid id specified", http.StatusBadRequest)
			return
		}

		sizeScannerId := SizeScannerId(id)
		delete(manager.SizeScanners, sizeScannerId)
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Invalid request method.", http.StatusNotImplemented)
	}

}

func (manager *ScannerManager) HandleRFIDRequest(w http.ResponseWriter, r *http.Request) {
	manager.RFIDScannerMutex.Lock()
	defer manager.RFIDScannerMutex.Unlock()

	switch r.Method {
	case "GET": //Register a new RFID scanner
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "No id specified", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid id specified", http.StatusBadRequest)
			return
		}

		manager.RFIDScanners = append(manager.RFIDScanners, RFIDScannerId(id))
		w.WriteHeader(http.StatusOK)
	case "POST": //Add a new RFID scan
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "No id specified", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid id specified", http.StatusBadRequest)
			return
		}

		valueStr := r.URL.Query().Get("value")
		if valueStr == "" {
			http.Error(w, "No value specified", http.StatusBadRequest)
			return
		}

		value, err := strconv.Atoi(valueStr)
		if err != nil {
			http.Error(w, "Invalid value specified", http.StatusBadRequest)
			return
		}

		manager.RFIDScanned <- RFIDScannerHistoryItem{RFIDScannerId(id), time.Now(), TrainId(value)}
		w.WriteHeader(http.StatusOK)

	case "DELETE": //Remove an RFID scanner
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "No id specified", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid id specified", http.StatusBadRequest)
			return
		}

		rfidScannerId := RFIDScannerId(id)
		for i, id := range manager.RFIDScanners {
			if id == rfidScannerId {
				manager.RFIDScanners = append(manager.RFIDScanners[:i], manager.RFIDScanners[i+1:]...)
				break
			}
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Invalid request method.", http.StatusNotImplemented)
	}

}
