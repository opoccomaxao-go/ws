package ws

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ws        *websocket.Conn
	errLogger *log.Logger
}

func NewClient(url string, errLogger *log.Logger) *Client {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		errLogger.Printf("WS Client.Dial %s\n", err.Error())
		panic(err)
	}
	return &Client{
		ws:        ws,
		errLogger: errLogger,
	}
}

func (c *Client) Pong() {
	if err := c.ws.WriteControl(websocket.PongMessage, nil, time.Now().Add(time.Minute)); err != nil {
		c.errLogger.Printf("WS Client.Pong %s\n", err.Error())
	}
}

func (c *Client) Ping() {
	if err := c.ws.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Minute)); err != nil {
		c.errLogger.Printf("WS Client.Ping %s\n", err.Error())
	}
}

func (c *Client) Read() (string, error) {
	_, msg, err := c.ws.ReadMessage()
	return string(msg), err
}

func (c *Client) Write(text string) error {
	return c.ws.WriteMessage(websocket.TextMessage, []byte(text))
}

func (c *Client) Close() {
	if err := c.ws.Close(); err != nil {
		c.errLogger.Printf("WS Client.Close %s\n", err.Error())
	}
}
