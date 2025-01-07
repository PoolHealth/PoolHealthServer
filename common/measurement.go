package common

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v5"
)

type Measurement struct {
	PoolID     uuid.UUID
	Chlorine   null.Float
	PH         null.Float
	Alkalinity null.Float
	CreatedAt  time.Time
}
