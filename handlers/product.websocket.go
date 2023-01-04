package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/GerardoHP/inventory-service/data"
	"golang.org/x/net/websocket"
)

type message struct {
	Data string `json: "data"`
	Type string `json: "type"`
}

func productSocket(ws *websocket.Conn) {
	go func(c *websocket.Conn) {
		for {
			var msg message
			if err := websocket.JSON.Receive(ws, &msg); err != nil {
				log.Println(err)
				break
			}

			fmt.Printf("received message %s \n", msg.Data)
		}
	}(ws)

	var repo data.ProductRepository = data.NewSqlRepository()
	for {
		products, err := repo.GetTopTenProducts()
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
