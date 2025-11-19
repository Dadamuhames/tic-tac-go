package handlers

import (
	"bytes"
	"log"
	"net/http"
	"xo-websocket/internal/clients"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func handleError(err error) {
	if err != nil {
		log.Println("error:", err)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true //r.Host == "xo-ui.vercel.app"
	},
}

var connectionPool = clients.NewConnectionPool()

func Handler(ctx *gin.Context) {
	w, r := ctx.Writer, ctx.Request
	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()

		if err != nil {
			log.Println("error:", err)
			break
		}

		if bytes.HasPrefix(message, []byte("register:")) {
			err := HandleRegistration(connectionPool, c, message)
			handleError(err)

		} else if bytes.HasPrefix(message, []byte("start:")) {
			err := HandleStartGame(connectionPool, c, message)
			handleError(err)

		} else {
			err := HandleMove(connectionPool, c, message)
			handleError(err)
		}

		log.Printf("recv:%s", message)
	}
}
