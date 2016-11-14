package model

import (
	"time"
)

// Config represents the configuration and is unique and holds data relevant for calculating the
// metabolic rate of the user
type Config struct {
	ID         int       `json:"id"`
	Height     float64   `json:"height"`
	Activity   float64   `json:"activity"`
	Birthday   time.Time `json:"birthday"`
	Gender     string    `json:"gender"`
	UnitSystem string    `json:"unitSystem"`
}
