package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/afterglowflexin/wildberries/level0/internal/cache"
	"github.com/afterglowflexin/wildberries/level0/internal/nats"
	"github.com/afterglowflexin/wildberries/level0/internal/server"
	"github.com/afterglowflexin/wildberries/level0/internal/store"
	"github.com/nats-io/stan.go"
)

func main() {
	//open DB connection
	//need no set db url
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	os.Setenv("DATABASE_URL", "postgres://1315474:44uthy9a@localhost:5432/wb")
	conn := store.OpenConnection()
	defer conn.Close(context.Background())

	//initialize cache
	cache := cache.New(conn)

	//connect to nats streaming
	sc := nats.NewStanServer()

	//subscribe to channel "orders" in nats streaming
	_, err := sc.Subscribe(
		"orders",
		func(msg *stan.Msg) {
			log.Printf("Found message %s", string(msg.Data))
			cache.AddOrder(string(msg.Data))
			store.AddOrder(msg.Data, conn)
		},
		stan.StartWithLastReceived(),
	)

	if err != nil {
		log.Fatal(err)
	}

	//start http server
	s := server.New(cache)
	s.Start()

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			log.Println("Received an interrupt, unsubscribing and closing connection...")

			sc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
