package main

import (
	a "an3softbot/internal/app"
	chcl "an3softbot/internal/transport"
	"context"
	"fmt"
	"log"
)

var App a.Application

func main() {
	c1 := chcl.ClickHouseClient{}
	conn, err := c1.Connect()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	rows, err := conn.Query(ctx, "SELECT Id, Message FROM Requests")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var (
			Id      uint64
			Message string
		)
		if err := rows.Scan(
			&Id,
			&Message,
		); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Id: %d, Message: \"%s\"", Id, Message)
	}

	App = a.Application{}
	//App.Run()
}
