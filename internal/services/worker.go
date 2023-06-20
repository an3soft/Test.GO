package services

import (
	t "an3softbot/internal/transport"
	"sync"
)

type Worker struct {
	TgBot *t.TelegramBot
	Wg    *sync.WaitGroup
}

func (worker *Worker) Run() {
	defer worker.Wg.Done()
}
