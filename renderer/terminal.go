package renderer

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/zupzup/calories/model"
	"github.com/zupzup/calories/util"
)

// TerminalRenderer is the renderer for the CLI
type TerminalRenderer struct{}

// Error renders an error
func (r *TerminalRenderer) Error(err error) (string, error) {
	return err.Error(), nil
}

// WeightHistory renders all weights in order and their dates
func (r *TerminalRenderer) WeightHistory(weights []model.Weight, config *model.Config) (string, error) {
	var res string
	for _, weight := range weights {
		res += fmt.Sprintf("\t%s: %s\n", weight.Created.Format(util.DateFormat), util.WeightUnit(config.UnitSystem, weight.Weight))
	}
	return fmt.Sprintf("Weight over time:\n%s\n", res), nil
}

// AddWeight renders a success message and the added weight
func (r *TerminalRenderer) AddWeight(weight float64, config *model.Config) (string, error) {
	return fmt.Sprintf("Set weight: %s \n", util.WeightUnit(config.UnitSystem, weight)), nil
}

// Config prints the given configuration with weight, amr and bmr
func (r *TerminalRenderer) Config(config *model.Config, weight *model.Weight, amr, bmr float64, age int) (string, error) {
	return fmt.Sprintf("Current Config:\n\tWeight: %s \n\tHeight: %s \n\tActivity: %.1f \n\tBirthday: %s (%d)\n\tGender: %s\n\tUnit System: %s\n\tAMR (BMR): %.0f (%.0f) calories per day\n",
		util.WeightUnit(config.UnitSystem, weight.Weight), util.HeightUnit(config.UnitSystem, config.Height), config.Activity, config.Birthday.Format(util.DateFormat), age, config.Gender, config.UnitSystem, amr, bmr), nil
}

// Days renders the days in the given timespan
func (r *TerminalRenderer) Days(days model.Days, from, to time.Time) (string, error) {
	res := fmt.Sprintf("Data from %s to %s:\n-----------------------------------\n", from.Format(util.DateFormat), to.Format(util.DateFormat))
	if len(days) > 0 {
		var formattedDays string
		sumAMR := 0.0
		sumCalories := 0
		for _, day := range days {
			sumAMR += getAMR(day)
			sumCalories += day.Used
			formattedDays += stringifyDay(day)
		}
		defSur := color.GreenString("deficit")
		result := sumAMR - float64(sumCalories)
		formattedResult := color.GreenString("%.0f", result)
		formattedCalories := color.GreenString("%d", sumCalories)
		if result < 0 {
			defSur = color.RedString("surplus")
			result = result * -1
			formattedResult = color.RedString("%.0f", result)
			formattedCalories = color.RedString("%d", sumCalories)
		}

		formattedDays += fmt.Sprintf("-----------------------------------\n%s / %.0f calories = %s %s\n", formattedCalories, sumAMR, formattedResult, defSur)
		return fmt.Sprintf("%s%s", res, formattedDays), nil
	}
	return fmt.Sprintf("%sNo entries have been found.\n", res), nil
}

// AddEntry displays a success message after adding an entry
func (r *TerminalRenderer) AddEntry(date string, calories int, food string) (string, error) {
	return fmt.Sprintf("Added Entry for %s with %d calories (%s)\n", date, calories, food), nil
}

// ClearEntries displays a success message after clearing the entries for a day
func (r *TerminalRenderer) ClearEntries(date string) (string, error) {
	return fmt.Sprintf("Cleared all entries for %s\n", date), nil
}

// ClearEntry displays a success message after clearing the entry at a given position for a day
func (r *TerminalRenderer) ClearEntry(date string, entry *model.Entry) (string, error) {
	return fmt.Sprintf("Cleared entry %d %s for %s\n", entry.Calories, entry.Food, date), nil
}

// stringifyDay turns a model.Day into it's terminal string representation
func stringifyDay(d *model.Day) string {
	var res string
	for i, entry := range d.Entries {
		if i == 0 {
			res += fmt.Sprintf("%s\n", entry.EntryDate)
		}
		res += fmt.Sprintf("\t%d %s\n", entry.Calories, entry.Food)
		if i == len(d.Entries)-1 {
			calorieString := color.GreenString("%d", d.Used)
			if float64(d.Used) > getAMR(d) {
				calorieString = color.RedString("%d", d.Used)
			}
			res += fmt.Sprintf("\t---------------------\n\t%s / %.0f calories\n", calorieString, getAMR(d))
		}
	}
	return res
}

// getAMR calculates the AMR for a whole model.Day
func getAMR(d *model.Day) float64 {
	if d.Entries != nil {
		return d.Entries[0].AMR
	}
	return 0
}

// Import displays a success message after importing from a file
func (r *TerminalRenderer) Import(fileName string, numEntries, numWeights int) (string, error) {
	return fmt.Sprintf("Imported data from %s with %d entries and %d weights\n", fileName, numEntries, numWeights), nil
}
