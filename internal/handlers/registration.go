package handlers

import (
	"encoding/json"
	"fmt"
	"strings"
	"xo-websocket/internal/clients"
	"xo-websocket/internal/dto"

	"github.com/gorilla/websocket"
)

func HandleRegistration(c *websocket.Conn, mt int, message []byte) error {
	name := strings.Replace(string(message), "register:", "", 1)
	newId := clients.AddClient(name, c)

	clientList := clients.GetAll()

	newUserMessage := make([]byte, 0)

	newUserMessage = fmt.Appendf(newUserMessage, `
				{"messageType": 0, "user": {"id": %d, "name": "%s"}}`, newId, name)

	for i := range clientList {
		client := clientList[i]

		if client.Id != newId {
			go client.Conn.WriteMessage(mt, newUserMessage)
		}
	}

	mappedClients := make([]dto.ClientResponse, 0)

	for i := range clientList {
		client := clientList[i]

		if client.Id == newId || client.Id == 0 {
			continue
		}

		mappedClients = append(mappedClients, dto.ClientResponse{
			Id:   client.Id,
			Name: client.Name,
		})
	}

	registerResponse := dto.RegisterResponse{
		MessageType: 3,
		User:        dto.ClientResponse{Id: newId, Name: name},
		Users:       mappedClients,
	}

	registerResponseJson, err := json.Marshal(registerResponse)

	if err != nil {
		return err
	}

	go c.WriteMessage(mt, registerResponseJson)

	return nil
}
