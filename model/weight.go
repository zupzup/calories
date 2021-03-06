package model

import (
	"time"
)

// Weight represents the user's weight at a given time
type Weight struct {
	ID      int       `storm:"id,increment" json:"id"`
	Created time.Time `json:"created"`
	Weight  float64   `json:"weight"`
}
