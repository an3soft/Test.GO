package transport

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ClickHouseClient struct {
}

func (cl *ClickHouseClient) Connect() (driver.Conn, error) {
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
			// Debugf: func(format string, v ...interface{}) {
			// 	fmt.Printf(format, v)
			// },
			// TLS: &tls.Config{
			// 	InsecureSkipVerify: true,
			// },
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

	v, err := conn.ServerVersion()
	if err != nil {
		return nil, err
	}
	println(v)

	return conn, nil
}
