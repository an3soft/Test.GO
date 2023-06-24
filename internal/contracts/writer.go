package contracts

import (
	m "an3softbot/internal/models"
	"context"
)

type DBWriter interface {
	Write(context.Context, *m.Request)
	Delete(context.Context, *m.Request)
}
