package handlers

import (
	"encoding/json"
	"fmt"
	"strings"
	"xo-websocket/internal/clients"
	"xo-websocket/internal/dto"

	"github.com/gorilla/websocket"
)

func HandleRegistration(cp *clients.ConnectionPool, c *websocket.Conn, message []byte) error {
	name := strings.Replace(string(message), "register:", "", 1)
	newId := cp.AddClient(name, c)

	c.SetCloseHandler(func(code int, text string) error {
		err := HandleCLose(cp, newId)
		return err
	})

	newUserMessage := make([]byte, 0)

	newUserMessage = fmt.Appendf(newUserMessage, `
				{"messageType": 0, "user": {"id": %d, "name": "%s"}}`, newId, name)

	cp.BroadcastMessage(newId, newUserMessage)

	mappedClients := make([]dto.ClientResponse, 0)

	clientList := cp.GetAll()

	for _, client := range clientList {
		if (*client).Id == newId || client.Id == 0 {
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

	go c.WriteMessage(websocket.TextMessage, registerResponseJson)

	return nil
}
