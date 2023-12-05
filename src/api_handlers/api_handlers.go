package apihandlers

import (
	"L0/src/db"
	"L0/src/models"
	"L0/src/nats"
	"encoding/json"
	"io"
	"net/http"
)

func GetOrderByID(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("orderID")

	rows, err := db.DB.Query("SELECT order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardKey, sm_id, date_created, oof_shard FROM orders WHERE order_uid = $1", orderID)
	if err != nil {
		http.Error(w, "Error while making query to DB", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		var deliveryJSON, paymentJSON, itemsJSON []byte

		err := rows.Scan(
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
			http.Error(w, "Error while scanning rows from DB", http.StatusInternalServerError)
			return
		}

		var delivery models.Delivery
		if err := json.Unmarshal(deliveryJSON, &delivery); err != nil {
			http.Error(w, "Error encoding delivery to JSON", http.StatusInternalServerError)
			return
		}
		order.Delivery = delivery

		var payment models.Payment
		if err := json.Unmarshal(paymentJSON, &payment); err != nil {
			http.Error(w, "Error encoding payment to JSON", http.StatusInternalServerError)
			return
		}
		order.Payment = payment

		var items []models.Item
		if err := json.Unmarshal(itemsJSON, &items); err != nil {
			http.Error(w, "Error encoding items to JSON", http.StatusInternalServerError)
			return
		}
		order.Items = items

		orderJSON, err := json.Marshal(order)
		if err != nil {
			http.Error(w, "Error encoding order to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(orderJSON)
	}

}

func PostMessageToNats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	nats.SendMessageToNats(w, string(body))
}
