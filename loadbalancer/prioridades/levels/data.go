package levels

import (
	"time"
)

type ServerInformation struct {
	ID        string
	Host      string
	Port      string
	Lambda    float64
	Timestamp time.Time
	Nclients  int
}
