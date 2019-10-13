package actions

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gorilla/websocket"
)

var validOrigins = []string{
	"moz-extension",
	"http://localhost",
	"http://127.0.0.1",
}

type WebSocketContext struct {
	buffalo.Context
	ws *websocket.Conn
}

func validateOriginFn(c buffalo.Context) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		origin := r.Header["Origin"][0]
		c.Logger().Printf("Validating Websocket origin: %s", origin)

		for _, o := range validOrigins {
			if strings.HasPrefix(origin, o) {
				return true
			}
		}

		return false
	}
}

func WebSocketMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     validateOriginFn(c),
		}
		c.Logger().Printf("Attempting to upgrade connection to WebSocket")
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)

		if err != nil {
			c.Logger().Errorf("Failed to upgrade connection: %v", err)
			err := struct {
				Success bool   `json:"success"`
				Message string `json:"message"`
			}{false, fmt.Sprintf("failed to upgrade connection to websocket: %v", err)}
			c.Render(500, r.JSON(err))
		}

		return next(WebSocketContext{c, conn})
	}

}

func WebsocketHandler(c buffalo.Context) error {
	ctx, ok := c.(WebSocketContext)
	if !ok {
		log.Fatalln("Context is not of type WebSocketContext")
	}
	for {
		msgType, data, err := ctx.ws.ReadMessage()
		logger := ctx.Logger()

		if err != nil {
			logger.Error(err)
			return nil
		}

		switch msgType {
		case websocket.BinaryMessage:
			{
				logger.Printf("Received Binary msg: %v", string(data))
				reply := struct {
					Type string
					URL  string
				}{
					"link",
					"https://www.youtube.com/watch?v=ce3Zfc7R6Zw",
				}

				msg, err := json.Marshal(reply)
				if err != nil {
					logger.Printf("Failed to marshall to json: %v", err)
				}

				logger.Printf("Sending binary message %d", len(msg))
				ctx.ws.WriteMessage(websocket.BinaryMessage, msg)
				time.AfterFunc(10*time.Second, func() {
					ctx.ws.WriteMessage(websocket.BinaryMessage, msg)
				})
			}
			break
		case websocket.TextMessage:
			{
				logger.Printf("Received Text msg: %v", string(data))
				v, err := strconv.Atoi(string(data))
				logger.Printf("Value: %v", v)
				if err != nil {
					logger.Printf("Failed to convert '%s' to int", data)
				} else {
					msg := strconv.Itoa(v + 1)
					logger.Printf("Sending: %v", msg)
					ctx.ws.WriteMessage(websocket.TextMessage, []byte(msg))
				}
			}
			break
		case websocket.CloseMessage:
			{
				logger.Printf("Websocket Closed")
			}
			break
		case websocket.PingMessage:
			fallthrough
		case websocket.PongMessage:
			fallthrough
		default:
			{
				logger.Printf("Received Msg Type: %v", string(msgType))
				logger.Printf("Data: %v", string(data))
			}
		}
	}
	return nil
}
