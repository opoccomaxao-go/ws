package ws

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	handler   Handler
	srv       *http.Server
	mux       *http.ServeMux
	errLogger *log.Logger
}

type Handler func(ws *websocket.Conn)

var upgrader = websocket.Upgrader{} // use default options

func NewServer(handler Handler, errLogger *log.Logger) *Server {
	res := Server{
		handler:   handler,
		errLogger: errLogger,
	}
	res.mux = http.NewServeMux()
	res.mux.Handle("/", http.FileServer(http.Dir("client")))
	res.mux.HandleFunc("/ws", res.handle)
	return &res
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.errLogger.Printf("WS Server.Upgrade %s\n", err.Error())
		return
	}
	s.handler(c)
	_ = c.Close()
}

func (s *Server) Listen(port string) {
	go func() {
		s.srv = &http.Server{
			Handler:      s.mux,
			Addr:         "127.0.0.1:" + port,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		err := s.srv.ListenAndServe()
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}()
}

func (s *Server) Stop() {
	_ = s.srv.Close()
}
