package command

import (
	"errors"
	"fmt"
	"github.com/zupzup/calories/datasource"
	"github.com/zupzup/calories/model"
	"github.com/zupzup/calories/renderer"
	"github.com/zupzup/calories/util"
	"math"
	"sort"
	"sync"
	"time"
)

// DayCommand is the command to show a range of days
type DayCommand struct {
	DataSource  datasource.DataSource
	Renderer    renderer.Renderer
	Week        bool
	Month       bool
	History     int
	DefaultDate string
}

// Execute shows the current day, if no parameters are used,
// otherwise shows the days for the given time span (day, week, month, history of days)
func (c *DayCommand) Execute() (string, error) {
	now := time.Now()
	fromDate := now
	toDate := now
	if c.Week {
		fromDate = util.GetBeginningOfWeek(now)
	} else if c.Month {
		fromDate = now.AddDate(0, 0, -now.Day()+1)
	} else if c.History != 0 {
		amount := c.History
		if amount < 0 {
			amount = amount * -1
		}
		fromDate = now.AddDate(0, 0, -amount)
	} else if c.DefaultDate != "" {
		parsedDate, err := time.Parse(util.DateFormat, c.DefaultDate)
		if err != nil {
			return "", fmt.Errorf("wrong format for date: %v, please use dd.mm.yyyy", err)
		}
		fromDate = parsedDate
		toDate = parsedDate
	}
	if toDate.Before(fromDate) {
		return "", errors.New("from-date needs to be before to-date")
	}
	days, err := fetchDuration(c.DataSource, fromDate, toDate)
	if err != nil {
		return "", err
	}
	return c.Renderer.Days(days, fromDate, toDate)
}

// fetchDuration fetches all days in the given timespan concurrently, with at most 20
// goroutines at the same time. After fetching, the list of days is sorted
// by date.
func fetchDuration(ds datasource.DataSource, from, to time.Time) (model.Days, error) {
	var wg sync.WaitGroup
	lock := sync.RWMutex{}
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)
	pool := make(chan struct{}, 20)
	diffInDays := int(math.Ceil(to.Sub(from).Hours()/24)) + 1

	days := model.Days{}
	for i := 0; i < diffInDays; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			pool <- struct{}{}
			defer func() { <-pool }()
			entryDate := from.AddDate(0, 0, i)
			formattedDate := entryDate.Format(util.DateFormat)
			entries, err := ds.FetchEntries(formattedDate)
			if err != nil {
				errChan <- err
			}
			if len(entries) > 0 {
				day := newDay(entries, entryDate)
				lock.Lock()
				defer lock.Unlock()
				days = append(days, day)
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		if err != nil {
			return nil, err
		}
	}
	sort.Sort(days)
	return days, nil
}

// newDay is the constructor for Day, calculates the used calories
func newDay(entries model.Entries, entryDate time.Time) *model.Day {
	used := 0
	for _, entry := range entries {
		used += entry.Calories
	}
	return &model.Day{Entries: entries, Used: used, Date: entryDate}
}
