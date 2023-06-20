package contracts

import (
	m "an3softbot/internal/models"
)

type Writer interface {
	Write(m.Request)
}
