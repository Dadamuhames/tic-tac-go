package handlers

import (
	"encoding/json"
	"fmt"
	"xo-websocket/internal/clients"
	"xo-websocket/internal/dto"
	"xo-websocket/internal/maps"

	"github.com/gorilla/websocket"
)

func HandleMove(c *websocket.Conn, mt int, message []byte) error {
	var moveMessage dto.MoveMessage

	err := json.Unmarshal(message, &moveMessage)

	if err != nil {
		return err
	}

	sendTo := clients.GetUserById(moveMessage.ToUserId)

	go clients.SendMessageToClient(sendTo, mt, message)

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

		c.WriteMessage(mt, winJson)

		winMessage.GameId = moveMessage.FromUserId

		winJson, err = json.Marshal(winMessage)

		if err != nil {
			return err
		}

		clients.SendMessageToClient(sendTo, mt, winJson)
	}

	return nil
}
