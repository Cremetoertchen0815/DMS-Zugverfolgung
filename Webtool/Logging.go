package webtool

import (
	"encoding/json"
	"time"
)

type LogItem struct {
	TimeStamp time.Time
	Type      string
	Message   string
}

type convertedLogItem struct {
	TimeStamp string `json:"timestamp"`
	Type      string `json:"type"`
	Message   string `json:"message"`
}

func (item *LogItem) String() (string, error) {
	convertedData := convertedLogItem{
		TimeStamp: item.TimeStamp.Format("2006-01-02 15:04:05"),
		Type:      item.Type,
		Message:   item.Message,
	}

	data, err := json.Marshal(convertedData)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
