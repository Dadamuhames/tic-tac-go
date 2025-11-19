package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"xo-websocket/internal/clients"
	"xo-websocket/internal/dto"
	"xo-websocket/internal/maps"

	"github.com/gorilla/websocket"
)

func HandleStartGame(c *websocket.Conn, mt int, message []byte) error {
	startMessage := bytes.Replace(message, []byte("start:"), []byte(""), 1)

	fmt.Println(string(startMessage))

	var start dto.StartRequest

	err := json.Unmarshal(startMessage, &start)

	if err != nil {
		return err
	}

	startWith := clients.GetUserById(start.ToUserId)

	startResponse := dto.StartResponse{
		MessageType: 1,
		FromUserId:  start.FromUserId,
		ToUserId:    start.ToUserId,
	}

	randomNumber := rand.IntN(2)

	nextXO := "X"

	if randomNumber == 0 {
		startResponse.Xo = "O"
	} else {
		startResponse.Xo = "X"
		nextXO = "O"
	}

	responseBytes, err := json.Marshal(startResponse)

	go clients.SendMessageToClient(startWith, mt, responseBytes)

	startResponse.Xo = nextXO

	responseBytes, err = json.Marshal(startResponse)

	c.WriteMessage(mt, responseBytes)

	maps.InitMap(start.FromUserId, start.ToUserId)

	return nil
}
