package app

import (
	c "an3softbot/internal/contracts"
	clkhs "an3softbot/internal/database/clickhouse"
	s "an3softbot/internal/services"
	t "an3softbot/internal/transport"
	"context"
	"os"
	"path/filepath"
	"sync"
)

type Application struct {
	AppPath string
	tgbot   *t.TelegramBot
	worker  *s.Worker
	writer  c.DBWriter
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
	ctx := context.Background()
	// Подкобчение к БД
	clkhs := clkhs.ClickHouseClient{ReadBufferSize: 3}
	_, err = clkhs.Connect(ctx)
	if err != nil {
		panic(err)
	}
	app.writer = &clkhs
	// Создание и запуск обработчика бота
	app.tgbot = &t.TelegramBot{TimeOut: 60,
		Debug:  false,
		Wg:     &app.wg,
		Writer: &clkhs}
	app.tgbot.Init()
	go app.tgbot.Run(ctx)
	// Создание и запуск обработчика запросов
	app.worker = &s.Worker{Reader: &clkhs,
		Writer:    &clkhs,
		Sender:    app.tgbot,
		Processor: &s.Processor{Writer: &clkhs},
		Wg:        &app.wg}
	go app.worker.Run(ctx)
	// Ожидание
	app.wg.Wait()
}
