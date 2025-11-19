package clients

import (
	"fmt"
	"log"
	"maps"
	"slices"

	"github.com/gorilla/websocket"
)

type Client struct {
	Id   int
	Name string
	Conn *websocket.Conn
}

var clients = make(map[int]*Client)

func AddClient(name string, conn *websocket.Conn) int {
	existsing := FindByName(name)

	if existsing != nil {
		existsing.Conn = conn
		fmt.Println("Connection replaced")
		return existsing.Id
	}

	id := len(clients) + 1

	clients[id] = &Client{Id: id, Name: name, Conn: conn}

	return id
}

func GetUserById(id int) *Client {
	return clients[id]
}

func SendMessageToClient(c *Client, mt int, message []byte) {
	err := c.Conn.WriteMessage(mt, message)

	if err != nil {
		log.Println("start game error: ", err)
	}
}

func GetAll() []*Client {
	return slices.Collect(maps.Values(clients))
}

func FindByName(name string) *Client {
	clients := GetAll()

	for _, client := range clients {
		if client.Name == name {
			return client
		}
	}
	return nil
}
