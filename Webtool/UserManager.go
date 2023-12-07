package webtool

import (
	"log"
	"math/rand"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type DataStreamManager struct {
	Users     map[int]*websocket.Conn
	DataInput chan LogItem
	mutex     sync.Mutex
	upgrader  *websocket.Upgrader
}

func CreateDataStreamManager() *DataStreamManager {
	data := DataStreamManager{
		Users:     make(map[int]*websocket.Conn),
		DataInput: make(chan LogItem),
		mutex:     sync.Mutex{},
		upgrader:  &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024},
	}

	go data.SendData()

	return &data
}

func (manager *DataStreamManager) AcceptNew(w http.ResponseWriter, r *http.Request) {
	//Upgrade connection to websocket
	conn, err := manager.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %s", err)
	}

	//Lock the user list
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	//Generate new random ID
	nuId := 0
	for {
		nuId = rand.Int()
		if _, exists := manager.Users[nuId]; !exists {
			break
		}
	}

	//Add user to list
	manager.Users[nuId] = conn
	log.Printf("[#%d]Admin instance started", nuId)

	//Wait for connection to close
	go func() {
		for {
			if _, _, err := conn.NextReader(); err != nil {
				//Connection has closed
				conn.Close()

				log.Printf("[#%d]Admin instance closed", nuId)

				//Lock the user list
				manager.mutex.Lock()
				defer manager.mutex.Unlock()

				//Remove user from list
				delete(manager.Users, nuId)

				break
			}
		}
	}()

}

func (manager *DataStreamManager) SendData() {
	for data := range manager.DataInput {
		//Get data as string
		dataStr, err := data.String()
		if err != nil {
			log.Printf("Error marshalling data: %s", err)
			continue
		}

		//Lock the user list
		manager.mutex.Lock()

		//Send data to all users
		for id, user := range manager.Users {
			err := user.WriteMessage(websocket.TextMessage, []byte(dataStr))
			if err != nil {
				log.Printf("[#%d]Error writing to user: %s", id, err)
			}

		}

		//Unlock the user list
		manager.mutex.Unlock()
	}
}
