package services

import (
	c "an3softbot/internal/contracts"
	m "an3softbot/internal/models"
	"context"
	"sync"
	"time"
)

type Worker struct {
	Reader    c.DBReader
	Writer    c.DBWriter
	Sender    c.Sender
	Processor c.Processor
	Wg        *sync.WaitGroup
}

func (worker *Worker) Run(ctx context.Context) {
	defer worker.Wg.Done()

	//for ; ; time.Sleep(12 * time.Hour) {
	for ; ; time.Sleep(5 * time.Minute) {
		println("Worker cycle start at", time.Now().String())
		select {
		case <-ctx.Done():
			return
		default:
			ch := worker.Reader.Read(ctx)
			var (
				req *m.Request
				ok  bool = true
			)
			for ok {
				select {
				case req, ok = <-ch:
					if ok {
						answ := worker.Processor.Process(ctx, req)
						if answ.Text != "" {
							worker.Sender.Send(ctx, &answ)
						}
					}
				case <-ctx.Done():
					return
				}
			}
		}
	}
}
