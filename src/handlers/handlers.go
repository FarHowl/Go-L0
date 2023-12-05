package handlers

import (
	"L0/src/db"
	"L0/src/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/xeipuuv/gojsonschema"
)

var ordersMap map[string]models.Order = make(map[string]models.Order)

func HandleJsonMessage(m *nats.Msg) {
	var order models.Order

	// Get the absolute JSON schema path on different machines
	relativePathToSchema := "src/models/model_schema.json"
	absolutePathToSchema, err := filepath.Abs(relativePathToSchema)
	if err != nil {
		log.Fatal(err)
	}
	absolutePathToSchema = strings.ReplaceAll(absolutePathToSchema, "\\", "/")

	// Validating the received JSON
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + absolutePathToSchema)
	documentLoader := gojsonschema.NewStringLoader(string(m.Data))
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		log.Fatal(err)
	} else if !result.Valid() {
		var responseString string
		for _, desc := range result.Errors() {
			responseString = fmt.Sprintf("%s ", desc)
		}
		log.Fatal(responseString)
	}

	// Decoding JSON to Order struct
	jsonErr := json.Unmarshal(m.Data, &order)
	if jsonErr != nil {
		log.Fatal(err)
	}

	// Executing the query to add new order into the DB
	stmt, err := db.DB.Prepare("INSERT INTO orders (order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardKey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Parse structures to JSONs to save them to DB
	delivery, err := json.Marshal(order.Delivery)
	if err != nil {
		log.Fatal(err)
	}
	payment, err := json.Marshal(order.Payment)
	if err != nil {
		log.Fatal(err)
	}
	items, err := json.Marshal(order.Items)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(order.OrderUID, order.TrackNumber, order.Entry, delivery, payment, items, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey, order.SMID, order.DateCreated, order.OOFShard)
	if err != nil {
		log.Fatal(err)
	}

	// Save order to cache
	ordersMap[order.OrderUID] = models.Order{
		OrderUID:          order.OrderUID,
		TrackNumber:       order.TrackNumber,
		Entry:             order.Entry,
		Delivery:          order.Delivery,
		Payment:           order.Payment,
		Items:             order.Items,
		Locale:            order.Locale,
		InternalSignature: order.InternalSignature,
		CustomerID:        order.CustomerID,
		DeliveryService:   order.DeliveryService,
		ShardKey:          order.ShardKey,
		SMID:              order.SMID,
		DateCreated:       order.DateCreated,
		OOFShard:          order.OOFShard,
	}

	fmt.Println("Successfully stored into the Database!")
}

// Retrieve orders from DB and save them to cache
func RetrieveOrdersFromDB() {
	orders, err := db.DB.Query("SELECT * FROM orders")
	if err != nil {
		log.Fatal(err)
	}
	defer orders.Close()

	var deliveryJSON sql.RawBytes
	var paymentJSON sql.RawBytes
	var itemsJSON sql.RawBytes

	for orders.Next() {
		var order models.Order
		err := orders.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&deliveryJSON,
			&paymentJSON,
			&itemsJSON,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SMID,
			&order.DateCreated,
			&order.OOFShard,
		)
		if err != nil {
			log.Fatal(err)
		}

		// Decode JSONs to structs
		var delivery models.Delivery
		if err := json.Unmarshal(deliveryJSON, &delivery); err != nil {
			log.Fatal(err)
		}
		order.Delivery = delivery

		var payment models.Payment
		if err := json.Unmarshal(paymentJSON, &payment); err != nil {
			log.Fatal(err)
		}
		order.Payment = payment

		var items []models.Item
		if err := json.Unmarshal(itemsJSON, &items); err != nil {
			log.Fatal(err)
		}
		order.Items = items

		ordersMap[order.OrderUID] = order
	}
}
