package handlers

import (
	"encoding/json"
	"xo-websocket/internal/clients"
	"xo-websocket/internal/dto"
)

func HandleCLose(cp *clients.ConnectionPool, userId int) error {
	cp.RemoveClient(userId)

	closeMessage := dto.CloseMessage{MessageType: 5, UserId: userId}
	closeMessageJson, err := json.Marshal(closeMessage)

	if err != nil {
		return err
	}

	cp.BroadcastMessage(0, closeMessageJson)

	return nil
}
