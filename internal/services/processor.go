package services

import (
	m "an3softbot/internal/models"
	"context"
	"fmt"
)

type Processor struct {
}

func (p *Processor) Process(ctx context.Context, req *m.Request) m.Answer {
	ret := m.Answer{}
	ret.ChatId = req.ChatId
	ret.Text = fmt.Sprintf("Обработка запроса: <b>%s</b>", req.Text)
	return ret
}
