package model

import (
	"time"
)

// Entry can be added and removes and hold the date they have been added,
// the date they have been added to, the used calories and the food which has been consumed.
// Also, for each entry, the metabolic rates are calculated, for later bookkeeping
type Entry struct {
	ID        int       `storm:"id,increment" json:"id"`
	Created   time.Time `json:"created"`
	EntryDate string    `json:"entryDate"`
	Calories  int       `json:"calories"`
	Food      string    `json:"food"`
	BMR       float64   `json:"bmr"`
	AMR       float64   `json:"amr"`
}

// Entries is a custom slice type for a list of entries
type Entries []Entry
