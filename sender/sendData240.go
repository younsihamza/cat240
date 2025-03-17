package sender

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"radar240/global"
)

func Sender() {
	route := gin.Default()
	route.GET("/radar240", func(c *gin.Context) {
		// fmt.Println("trying to connect ")
		conn , err := global.Upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("ERROR: ", err)
			return
		}
		global.Clients[conn] = true
		defer func () {
			delete(global.Clients, conn)
			conn.Close()
		}()
		for {
			select {
				case message := <- global.ParsedData:
					for client := range global.Clients {
						err := client.WriteMessage(websocket.TextMessage, message)
						if err != nil {
							fmt.Println("ERROR: ", err)
							client.Close()
							delete(global.Clients, client)
						}
					}
				}
			}
	})
	route.Run(":8083")
}