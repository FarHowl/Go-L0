package main

import (
	apihandlers "L0/src/api_handlers"
	"L0/src/db"
	"L0/src/handlers"
	"L0/src/nats"
	"fmt"
	"net/http"
)

func main() {
	db.ConnectToDB()
	handlers.RetrieveOrdersFromDB()
	nats.StartNatsServer()

	http.HandleFunc("/getOrderByID", apihandlers.GetOrderByID)

	http.HandleFunc("/postMessageToNats", apihandlers.PostMessageToNats)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Print("Error wile starting the HTTP server")
	}
}
