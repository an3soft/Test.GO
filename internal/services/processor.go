package services

import (
	c "an3softbot/internal/contracts"
	m "an3softbot/internal/models"
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Processor struct {
	Writer c.DBWriter
}

func (p *Processor) Process(ctx context.Context, req *m.Request) m.Answer {
	ret := m.Answer{}
	ret.ChatId = req.ChatId
	// Generate random process result for stub
	rnd := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	var ready bool = rnd.Intn(2) == 1
	if ready {
		ret.Text = fmt.Sprintf("Запрос обработан: <b>%s</b>", req.Text)
	} else {
		ret.Text = fmt.Sprintf("Запрос НЕ обработан: <b>%s</b>", req.Text)
	}
	req.Ready = ready
	p.Writer.Write(ctx, req)
	return ret
}
