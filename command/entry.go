package command

import (
	"fmt"
	"github.com/zupzup/calories/datasource"
	"github.com/zupzup/calories/renderer"
	"github.com/zupzup/calories/util"
	"os"
	"strconv"
	"time"
)

// ClearEntriesCommand is the command to clear entries for a given day
type ClearEntriesCommand struct {
	DataSource datasource.DataSource
	Renderer   renderer.Renderer
	Date       string
	Position   int
	YesMode    bool
}

// Execute removes all entries of the current day, if no parameters are given
// otherwise it removes the entries of the given date.
// Asks the user for confirmation
func (c *ClearEntriesCommand) Execute() (string, error) {
	chosenDate := time.Now()
	if c.Date != "" {
		parsedDate, err := time.Parse(util.DateFormat, c.Date)
		if err != nil {
			return "", fmt.Errorf("wrong format for date: %v, please use dd.mm.yyyy", err)
		}
		chosenDate = parsedDate
	}
	formattedDate := chosenDate.Format(util.DateFormat)
	if c.Position >= 0 {
		return clearSingleEntry(c.DataSource, c.Renderer, c.YesMode, formattedDate, c.Position)
	}
	return clearAllEntries(c.DataSource, c.Renderer, c.YesMode, formattedDate)
}

// clearSingleEntry deletes a single entry based on the given position from the database after asking the user, validating the given position
func clearSingleEntry(ds datasource.DataSource, r renderer.Renderer, yesMode bool, formattedDate string, position int) (string, error) {
	entries, err := ds.FetchEntries(formattedDate)
	if err != nil {
		return "", err
	}
	if len(entries) == 0 {
		return "", fmt.Errorf("could not delete entry at position %d for %s, there are no entries", position, formattedDate)
	}
	if position == 0 || position > len(entries) {
		return "", fmt.Errorf("could not delete entry at position %d for %s, value needs to be from %d to %d", position, formattedDate, 1, len(entries))
	}
	entry := entries[position-1]
	if !yesMode {
		choice, confErr := util.AskConfirmation(fmt.Sprintf("Do you really want to clear the entry %d %s for %s? The data will be lost.", entry.Calories, entry.Food, formattedDate), os.Stdin)
		if confErr != nil {
			return "", confErr
		}
		if !choice {
			return "", nil
		}
	}
	err = ds.RemoveEntry(formattedDate, entry.ID)
	if err != nil {
		return "", err
	}
	return r.ClearEntry(formattedDate, &entry)
}

// clearAllEntries deletes all entries for a given day, after asking the user
func clearAllEntries(ds datasource.DataSource, r renderer.Renderer, yesMode bool, formattedDate string) (string, error) {
	if !yesMode {
		choice, err := util.AskConfirmation(fmt.Sprintf("Do you really want to clear all entries for %s? The data will be lost.", formattedDate), os.Stdin)
		if err != nil {
			return "", err
		}
		if !choice {
			return "", nil
		}
	}
	err := ds.RemoveEntries(formattedDate)
	if err != nil {
		return "", err
	}
	return r.ClearEntries(formattedDate)
}

// AddEntryCommand is the command to add an entry for a day
type AddEntryCommand struct {
	DataSource datasource.DataSource
	Renderer   renderer.Renderer
	Date       string
	Food       string
	Calories   string
	Mode       int
}

// Execute shows, if there is no date parameter given, the given calories and food are added to the current day,
// otherwise to the given date
func (c *AddEntryCommand) Execute() (string, error) {
	if c.Mode < 2 {
		return "", fmt.Errorf("usage: calories add [--d=DATE] [--o=FORMAT] CALORIES FOOD")
	}
	chosenDate := time.Now()
	calories, err := strconv.Atoi(c.Calories)
	if err != nil {
		return "", fmt.Errorf("wrong format for calories: %s needs to be a number (e.g.: 600)", c.Calories)
	}
	if c.Date != "" {
		parsedDate, parseErr := time.Parse(util.DateFormat, c.Date)
		if parseErr != nil {
			return "", fmt.Errorf("wrong format for date: %v, please use dd.mm.yyyy", parseErr)
		}
		chosenDate = parsedDate
	}
	formattedDate := chosenDate.Format(util.DateFormat)
	err = c.DataSource.AddEntry(formattedDate, calories, c.Food)
	if err != nil {
		return "", err
	}
	return c.Renderer.AddEntry(formattedDate, calories, c.Food)
}
