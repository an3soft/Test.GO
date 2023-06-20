package contracts

import (
	m "an3softbot/internal/models"
)

type Sender interface {
	Send(m.Answer)
}
