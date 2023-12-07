package main

import (
	"TestApp/api"
	"TestApp/webtool"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
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

	//Setup static file server
	fileServer := http.FileServer(http.Dir("Webtool/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	//Create webtool
	pageman := webtool.CreateWebtoolHandler()
	userman := webtool.CreateDataStreamManager()
	http.HandleFunc("/admin", pageman.HandleMainPage)
	http.HandleFunc("/ws", userman.AcceptNew)

	//Create a new scanner manager
	manager := api.CreateScannerManager(userman.DataInput)
	http.HandleFunc("/api/scanner/size", manager.HandleSizeRequest)
	http.HandleFunc("/api/scanner/rfid", manager.HandleRFIDRequest)
	http.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) { ProcessPingRequest(w, r, userman.DataInput) })

	//Start data processing threads
	go ProcessSizeData(manager.SizeScanned)
	go ProcessRFIDData(manager.RFIDScanned)

	http.ListenAndServe("localhost:8080", nil)
	fmt.Println("Starting server...")
}

func ProcessSizeData(input chan api.SizeScannerHistoryItem) {
	for data := range input {
		//TODO: Update train size in database
		_ = data
	}
}

func ProcessRFIDData(dataFunnel chan api.RFIDScannerHistoryItem) {
	for data := range dataFunnel {
		//TODO: Update train size in database
		_ = data
	}

}

func ProcessPingRequest(w http.ResponseWriter, r *http.Request, webOutput chan<- webtool.LogItem) {
	//Send to active websocket connections
	webOutput <- webtool.LogItem{
		TimeStamp: time.Now(),
		Type:      "ADMIN",
		Message:   "Ping received by server!"}

	w.WriteHeader(http.StatusOK)
}
