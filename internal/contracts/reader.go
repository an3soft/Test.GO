package contracts

import (
	m "an3softbot/internal/models"
	"context"
)

type DBReader interface {
	Read(context.Context) chan *m.Request
}
