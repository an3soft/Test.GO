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
	// Generating a random processing result for a stub
	rnd := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	var ready bool = rnd.Intn(2) == 1
	if ready {
		ret.Text = fmt.Sprintf("Запрос обработан: <b>%s</b>", req.Text)
		//req.Ready = true
		//p.Writer.Write(ctx, req)	// clickhouse merge is bad, if you need update, use specific version field for select last update (or secondary table)
		p.Writer.Delete(ctx, req) // delete work best
	} else {
		ret.Text = fmt.Sprintf("Запрос НЕ обработан: <b>%s</b>", req.Text)
		//p.Writer.Write(ctx, req)	// clickhouse merge is bad
	}
	return ret
}
