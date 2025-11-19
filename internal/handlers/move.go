package handlers

import (
	"encoding/json"
	"fmt"
	"xo-websocket/internal/clients"
	"xo-websocket/internal/dto"
	"xo-websocket/internal/maps"

	"github.com/gorilla/websocket"
)

func HandleMove(cp *clients.ConnectionPool, c *websocket.Conn, message []byte) error {
	var moveMessage dto.MoveMessage

	err := json.Unmarshal(message, &moveMessage)

	if err != nil {
		return err
	}

	sendTo := cp.GetClientById(moveMessage.ToUserId)

	go sendTo.SendMessage(message)

	maps.AddToMap(
		moveMessage.FromUserId,
		moveMessage.ToUserId,
		moveMessage.Row,
		moveMessage.Col,
		moveMessage.Xo)

	winner := maps.CheckWinner(moveMessage.FromUserId, moveMessage.ToUserId)

	fmt.Printf("Winner: %s\n", winner)

	if winner != "" {
		winMessage := dto.WinnerMessage{MessageType: 4, Winner: winner, GameId: moveMessage.ToUserId}

		winJson, err := json.Marshal(winMessage)

		if err != nil {
			return err
		}

		c.WriteMessage(websocket.TextMessage, winJson)

		winMessage.GameId = moveMessage.FromUserId

		winJson, err = json.Marshal(winMessage)

		if err != nil {
			return err
		}

		sendTo.SendMessage(winJson)
	}

	return nil
}
