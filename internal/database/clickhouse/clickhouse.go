package database

import (
	m "an3softbot/internal/models"
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ClickHouseClient struct {
	Connected  bool
	connection *driver.Conn
}

func (cl *ClickHouseClient) Connect() (driver.Conn, error) {
	if cl.Connected {
		return *cl.connection, nil
	}

	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{"localhost:19000"},
			Auth: clickhouse.Auth{
				Database: "an3softBotData",
				Username: "default",
				Password: "",
			},
			ClientInfo: clickhouse.ClientInfo{
				Products: []struct {
					Name    string
					Version string
				}{
					{Name: "an3softBot", Version: "0.1"},
				},
			},
		})
	)

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}

	cl.connection = &conn
	cl.Connected = true

	return conn, nil
}

func (cl *ClickHouseClient) Write(request m.Request) {

	// v, err := cl.connection.ServerVersion()
	// if err != nil {
	// 	return nil, err
	// }
	// println(v)

	println("Write request:")
	println(request.UserId)
	println(request.ChatId)
	println(request.MessageID)
	println(request.UserName)
	println(request.Text)
	println(request.Received.String())
}
