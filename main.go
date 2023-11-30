package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	//Create log files
	fmt.Println("Open log...")
	logFile, err := os.OpenFile("./history.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	log.Println("---SERVER STARTED---")

	fmt.Println("Starting server...")

	//Create a new scanner manager
	manager := CreateScannerManager()
	http.HandleFunc("/api/scanner/size", manager.HandleSizeRequest)
	http.HandleFunc("/api/scanner/rfid", manager.HandleRFIDRequest)

	//Start data processing threads
	go ProcessSizeData(manager.SizeScanned)
	go ProcessRFIDData(manager.RFIDScanned)

	http.ListenAndServe("localhost:8080", nil)
	fmt.Println("Starting server...")
}

func ProcessSizeData(dataFunnel chan SizeScannerHistoryItem) {
	for data := range dataFunnel {
		//Log the data
		logStr := fmt.Sprintf("RFID scanned { Scanner: %d, Timetamp: %d-%d-%d %d:%d:%d.%d, Train: %d, Size: %f }",
			data.Scanner,
			data.TimeStamp.Year(),
			data.TimeStamp.Month(),
			data.TimeStamp.Day(),
			data.TimeStamp.Hour(),
			data.TimeStamp.Minute(),
			data.TimeStamp.Second(),
			data.TimeStamp.Nanosecond(),
			data.Train,
			data.Value)

		fmt.Println(logStr)
		log.Println(logStr)

		//TODO: Send to active websocket connections

		//TODO: Update train size in database
	}
}

func ProcessRFIDData(dataFunnel chan RFIDScannerHistoryItem) {
	for data := range dataFunnel {
		logStr := fmt.Sprintf("RFID scanned { Scanner: %d, Timetamp: %d-%d-%d %d:%d:%d.%d, RFID: %d }",
			data.Scanner,
			data.TimeStamp.Year(),
			data.TimeStamp.Month(),
			data.TimeStamp.Day(),
			data.TimeStamp.Hour(),
			data.TimeStamp.Minute(),
			data.TimeStamp.Second(),
			data.TimeStamp.Nanosecond(),
			data.Value)

		fmt.Println(logStr)
		log.Println(logStr)

		//TODO: Send to active websocket connections
		//TODO: Update train size in database
	}

}
