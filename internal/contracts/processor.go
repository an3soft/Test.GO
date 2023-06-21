package contracts

import (
	m "an3softbot/internal/models"
	"context"
)

type Processor interface {
	Process(context.Context, *m.Request) m.Answer
}
