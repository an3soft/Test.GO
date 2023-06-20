package app

import (
	c "an3softbot/internal/contracts"
	clkhs "an3softbot/internal/database/clickhouse"
	s "an3softbot/internal/services"
	t "an3softbot/internal/transport"
	"os"
	"path/filepath"
	"sync"
)

type Application struct {
	AppPath string
	tgbot   *t.TelegramBot
	worker  *s.Worker
	writer  c.Writer
	wg      sync.WaitGroup
}

func (app *Application) Run() {
	appPath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	app.AppPath, err = filepath.Abs(filepath.Dir(appPath))
	if err != nil {
		panic(err)
	}
	app.wg.Add(2)
	clkhs := clkhs.ClickHouseClient{}
	_, err = clkhs.Connect()
	if err != nil {
		panic(err)
	}
	app.writer = &clkhs
	// ctx := context.Background()
	// rows, err := conn.Query(ctx, "SELECT Id, Message FROM Requests")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for rows.Next() {
	// 	var (
	// 		Id      uint64
	// 		Message string
	// 	)
	// 	if err := rows.Scan(
	// 		&Id,
	// 		&Message,
	// 	); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("Id: %d, Message: \"%s\"", Id, Message)
	// }

	// Создание и запуск обработчика бота
	app.tgbot = &t.TelegramBot{TimeOut: 30, Debug: false, Wg: &app.wg, Writer: &clkhs}
	go app.tgbot.Run()
	// Создание и запуск обработчика запросов
	app.worker = &s.Worker{TgBot: app.tgbot, Wg: &app.wg}
	go app.worker.Run()
	// Ожидание
	app.wg.Wait()
}
