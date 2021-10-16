package ws

import (
	"io"
	"log"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func Echo(ws *websocket.Conn) {
	for {
		t, b, _ := ws.ReadMessage()
		_ = ws.WriteMessage(t, b)
	}
}

func TestServer(t *testing.T) {
	const port = "65534"
	s := NewServer(Echo, log.New(io.Discard, "", 0))
	s.Listen(port)
	wsClient := NewClient("ws://localhost:"+port+"/ws", log.New(io.Discard, "", 0))
	testStr := []string{"test", "test2", "test4"}
	for _, value := range testStr {
		_ = wsClient.Write(value)
		resp, _ := wsClient.Read()
		assert.Equal(t, value, resp, "Message")
	}
}
