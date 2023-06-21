package services

import (
	c "an3softbot/internal/contracts"
	"time"

	//m "an3softbot/internal/models"
	"context"
	"sync"
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
	for ; ; time.Sleep(1 * time.Minute) {
		select {
		case <-ctx.Done():
			return
		default:
			ch := worker.Reader.Read(ctx)
			for {
				select {
				case req := <-ch:
					answ := worker.Processor.Process(ctx, req)
					if answ.Text != "" {
						worker.Sender.Send(ctx, &answ)
					}
				case <-ctx.Done():
					return
				}
			}
		}
	}
}
