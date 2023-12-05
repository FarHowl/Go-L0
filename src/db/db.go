package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectToDB() {
	connectStr := "postgres://postgres:3498569@localhost/postgres?sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connectStr)
	if err != nil {
		log.Fatal(err)
	}

	// Create table `orders` if it does not exist
	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS orders (
            order_uid VARCHAR(255) PRIMARY KEY,
            track_number VARCHAR(255),
            entry VARCHAR(255),
            delivery JSON,
            payment JSON,
            items JSON,
            locale VARCHAR(10),
            internal_signature VARCHAR(255),
            customer_id VARCHAR(255),
            delivery_service VARCHAR(255),
            shardKey VARCHAR(255),
            sm_id INTEGER,
            date_created TIMESTAMP,
            oof_shard VARCHAR(255)
        )
    `)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Database")
}
