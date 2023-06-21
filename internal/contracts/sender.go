package contracts

import (
	m "an3softbot/internal/models"
	"context"
)

type Sender interface {
	Send(context.Context, *m.Answer)
}
