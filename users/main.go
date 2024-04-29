package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/dubey22rohit/heyyy_yo_backend/types"
	"github.com/gorilla/websocket"
)

const wsEndpoint = "ws://127.0.0.1:30000/ws"

var sendInterval = time.Second * 10

func generateUsername() string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers := "0123456789"

	lengthLetter := 5
	startLetter := rand.Intn(len(letters) - lengthLetter + 1)
	windowLetter := letters[startLetter : startLetter+lengthLetter]

	lengthNumber := 4
	startNumber := rand.Intn(len(numbers) - lengthNumber + 1)
	windowNumber := numbers[startNumber : startNumber+lengthNumber]

	username := windowLetter + windowNumber

	return username
}

func generateUser() types.User {
	user_id := rand.Intn(999999)
	username := generateUsername()
	var online bool
	if rand.Intn(2) == 0 {
		online = true
	} else {
		online = false
	}
	user := types.User{
		User_id:  int64(user_id),
		Username: username,
		Online:   online,
	}
	return user
}

func main() {
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		for i := 0; i < 10; i++ {
			user := generateUser()
			if err := conn.WriteJSON(user); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(sendInterval)
	}

}
