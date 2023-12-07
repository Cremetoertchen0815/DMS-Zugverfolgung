package api

import (
	"TestApp/webtool"
	"fmt"
	"log"
	"time"
)

func (data RFIDScannerHistoryItem) Log(webOutput chan<- webtool.LogItem) {
	//Log the data
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

	log.Println(logStr)

	//Send to active websocket connections
	webOutput <- webtool.LogItem{
		TimeStamp: data.TimeStamp,
		Type:      "RFID",
		Message:   fmt.Sprintf("Train scanned! (ID: %d)", data.Value)}
}

func (data SizeScannerHistoryItem) Log(webOutput chan<- webtool.LogItem) {
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

	//Send to active websocket connections
	webOutput <- webtool.LogItem{
		TimeStamp: data.TimeStamp,
		Type:      "SIZE",
		Message:   fmt.Sprintf("Train scanned! (ID: %d, Size: %f)", data.Train, data.Value)}
}

func (manager *ScannerManager) LogRFIDMessage(message string) {
	//Log the data
	log.Println(message)

	//Send to active websocket connections
	manager.Logger <- webtool.LogItem{
		TimeStamp: time.Now(),
		Type:      "RFID",
		Message:   message}
}

func (manager *ScannerManager) LogSizeMessage(message string) {
	//Log the data
	log.Println(message)

	//Send to active websocket connections
	manager.Logger <- webtool.LogItem{
		TimeStamp: time.Now(),
		Type:      "SIZE",
		Message:   message}
}
