package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/GerardoHP/inventory-service/data"
	"github.com/GerardoHP/inventory-service/models"
	"golang.org/x/net/websocket"
)

func ProductSocket(ws *websocket.Conn) {
	done := make(chan struct{})
	go func(c *websocket.Conn) {
		for {
			var msg models.Message
			if err := websocket.JSON.Receive(ws, &msg); err != nil {
				log.Println(err)
				break
			}

			fmt.Printf("received message %s \n", msg.Data)
		}
		close(done)
	}(ws)

	var repo data.ProductRepository = data.NewSqlRepository()
loop:
	for {
		select {
		case <-done:
			fmt.Println("connection was closed, lets break out of here")
			break loop
		default:
			products, err := repo.GetTopProducts(10)
			if err != nil {
				log.Println(err)
				break
			}

			if err := websocket.JSON.Send(ws, products); err != nil {
				log.Println(err)
				break
			}

			time.Sleep(10 * time.Second)
		}
	}

	fmt.Println("closing the websocket")
	defer ws.Close()
}
