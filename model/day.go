package model

import (
	"time"
)

// Day is an actual day with all it's entries and the
// calories which have been used for the day
type Day struct {
	Entries Entries   `json:"entries"`
	Used    int       `json:"used"`
	Date    time.Time `json:"date"`
}

// Days is Custom slice type for a list of days
type Days []*Day

func (days Days) Len() int           { return len(days) }
func (days Days) Less(i, j int) bool { return days[i].Date.Before(days[j].Date) }
func (days Days) Swap(i, j int)      { days[i], days[j] = days[j], days[i] }
