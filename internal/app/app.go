package app

import (
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
	// Создание и запуск обработчика бота
	app.tgbot = &t.TelegramBot{TimeOut: 30, Debug: false, Wg: &app.wg}
	go app.tgbot.Run()
	// Создание и запуск обработчика запросов
	app.worker = &s.Worker{TgBot: app.tgbot, Wg: &app.wg}
	go app.worker.Run()
	// Ожидание
	app.wg.Wait()
}
