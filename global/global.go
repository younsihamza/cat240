package global

import (
	"net/http"

	"github.com/gorilla/websocket"
)


var (
	FilteredData = make(chan []byte, 30)
	ParsedData = make(chan []byte, 30)
	Clients = make(map[*websocket.Conn]bool)

)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},}