package nats

import (
	"L0/src/handlers"
	"log"
	"net/http"

	"github.com/nats-io/nats.go"
)

var Nc *nats.Conn
var subject = "main"

func StartNatsServer() {
	// Connection to NATS Streaming Server
	var err error
	Nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	// Subscribe to the NATS broker
	Nc.Subscribe(subject, func(m *nats.Msg) {
		handlers.HandleJsonMessage(m)
	})
}

// We use transactions in order to have a response from NATS that is taken to channel
func SendMessageToNats(w http.ResponseWriter, jsonData string) {
	// transactionID := transactions.GenerateTransactionID()
	// ch := make(chan string, 1)

	// if transactions.Transactions == nil {
	// 	transactions.Transactions = make(map[string]chan string)
	// }
	// transactions.Transactions[transactionID] = ch

	// var order models.Order
	// err := json.Unmarshal([]byte(jsonData), &order)
	// if err != nil {
	// 	http.Error(w, "Error while parsing string to JSON", http.StatusInternalServerError)
	// 	return
	// }

	// order.ChannelTransactionId = transactionID
	// newJson, err := json.Marshal(order)
	// if err != nil {
	// 	http.Error(w, "Error while parsing JSON to string", http.StatusInternalServerError)
	// 	return
	// }

	if err := Nc.Publish(subject, []byte(jsonData)); err != nil {
		http.Error(w, "Error while publishing the message to NATS Streaming", http.StatusInternalServerError)
		return
	}

	// fmt.Println(transactions.Transactions)

	// result := <-ch
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte(result))
	// close(ch)
	// delete(transactions.Transactions, transactionID)
}
