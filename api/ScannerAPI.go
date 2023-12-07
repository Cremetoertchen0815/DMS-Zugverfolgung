package api

import (
	"TestApp/bl"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

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
		manager.LogSizeMessage(fmt.Sprintf("Scanner #%d added!", id))

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

		result := SizeScannerHistoryItem{SizeScannerId(id), time.Now(), trainId, value}
		result.Log(manager.Logger)
		manager.SizeScanned <- result

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

		manager.LogSizeMessage(fmt.Sprintf("Scanner #%d removed!", sizeScannerId))

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
		manager.LogRFIDMessage(fmt.Sprintf("Scanner #%d added!", id))

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

		result := RFIDScannerHistoryItem{RFIDScannerId(id), time.Now(), bl.TrainId(value)}

		manager.RFIDScanned <- result
		result.Log(manager.Logger)

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

		manager.LogRFIDMessage(fmt.Sprintf("Scanner #%d removed!", rfidScannerId))

		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Invalid request method.", http.StatusNotImplemented)
	}

}
