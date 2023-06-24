package database

import (
	m "an3softbot/internal/models"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ClickHouseClient struct {
	Connected      bool
	connection     *driver.Conn
	ReadBufferSize int
}

func (cl *ClickHouseClient) Connect(ctx context.Context) (driver.Conn, error) {
	if cl.Connected {
		return *cl.connection, nil
	}

	var (
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

func GetRow(rows driver.Rows) (*m.Request, error) {
	if rows.Next() {
		req := m.Request{}
		var mId int32
		if err := rows.Scan(
			&req.UserId,
			&req.ChatId,
			&mId,
			&req.UserName,
			&req.Text,
			&req.Received,
			&req.Updated,
			&req.Ready,
		); err != nil {
			return nil, err
		}
		req.MessageID = int(mId)
		//fmt.Printf("ChatId: %d, Message: \"%s\"\n\r", req.ChatId, req.Text)
		return &req, nil
	}

	return nil, nil
}

func (cl *ClickHouseClient) Write(ctx context.Context, request *m.Request) {

	conn := *cl.connection
	// Requests ENGINE = ReplacingMergeTree
	rows, err := conn.Query(ctx, fmt.Sprintf("SELECT * FROM Requests FINAL WHERE ChatId = %d AND MessageID = %d", request.ChatId, request.MessageID))
	if err != nil {
		log.Fatal(err)
	}
	req, err := GetRow(rows)
	if err != nil {
		log.Fatal(err)
	}

	if req == nil {
		err := conn.Exec(ctx, fmt.Sprintf("INSERT INTO Requests VALUES(%d, %d, %d, '%s', '%s', '%s', '%s', false)",
			request.UserId,
			request.ChatId,
			request.MessageID,
			request.UserName,
			request.Text,
			time.Now().Format("2006-01-02 15:04:05"),
			time.Now().Format("2006-01-02 15:04:05")))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Requests ENGINE = ReplacingMergeTree
		err := conn.Exec(ctx, fmt.Sprintf("INSERT INTO Requests VALUES(%d, %d, %d, '%s', '%s', '%s', '%s', %t)",
			request.UserId,
			request.ChatId,
			request.MessageID,
			request.UserName,
			request.Text,
			request.Received.Format("2006-01-02 15:04:05"),
			time.Now().Format("2006-01-02 15:04:05"),
			request.Ready))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (cl *ClickHouseClient) Read(ctx context.Context) chan *m.Request {
	conn := *cl.connection
	rows, err := conn.Query(ctx, fmt.Sprintf("SELECT * FROM Requests WHERE Ready = %t", false))
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan *m.Request, cl.ReadBufferSize)

	go func() {
		req, err := GetRow(rows)
		if err != nil {
			log.Fatal(err)
		}
		for req != nil {
			select {
			case <-ctx.Done():
				close(ch)
				return
			default:
				ch <- req
				req, err = GetRow(rows)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		close(ch)
	}()

	return ch
}

func (cl *ClickHouseClient) Delete(ctx context.Context, request *m.Request) {

	conn := *cl.connection
	// Requests ENGINE = ReplacingMergeTree
	rows, err := conn.Query(ctx, fmt.Sprintf("SELECT * FROM Requests FINAL WHERE ChatId = %d AND MessageID = %d", request.ChatId, request.MessageID))
	if err != nil {
		log.Fatal(err)
	}
	req, err := GetRow(rows)
	if err != nil {
		log.Fatal(err)
	}

	if req != nil {
		err := conn.Exec(ctx, fmt.Sprintf("DELETE FROM Requests WHERE ChatId = %d AND MessageID = %d", request.ChatId, request.MessageID))
		if err != nil {
			log.Fatal(err)
		}
	}
}
