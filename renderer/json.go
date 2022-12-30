package renderer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/zupzup/calories/model"
	"github.com/zupzup/calories/util"
)

type success struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// JSONRenderer is the JSON renderer
type JSONRenderer struct{}

// Error renders an error
func (r *JSONRenderer) Error(err error) (string, error) {
	res := success{
		Success: false,
		Message: err.Error(),
	}
	b, err := json.Marshal(res)
	if err != nil {
		return "", fmt.Errorf("could not marshal json, %v", err)
	}
	return string(b), nil
}

// WeightHistory renders all weights in order and their dates
func (r *JSONRenderer) WeightHistory(weights []model.Weight, config *model.Config) (string, error) {
	type weightUnit struct {
		Created   time.Time `json:"created"`
		Weight    float64   `json:"weight"`
		Formatted string    `json:"formatted"`
	}
	var res []*weightUnit
	for _, w := range weights {
		weight := w.Weight
		if config.UnitSystem == util.Imperial {
			weight = util.ToPounds(weight)
		}
		res = append(res, &weightUnit{
			Created:   w.Created,
			Weight:    weight,
			Formatted: util.WeightUnit(config.UnitSystem, w.Weight),
		})
	}
	b, err := json.Marshal(res)
	if err != nil {
		return "", fmt.Errorf("could not marshal json, %v", err)
	}
	return string(b), nil
}

// AddWeight renders a success message and the added weight
func (r *JSONRenderer) AddWeight(weight float64, config *model.Config) (string, error) {
	res := success{
		Success: true,
		Message: fmt.Sprintf("Set weight: %s", util.WeightUnit(config.UnitSystem, weight)),
	}
	b, err := json.Marshal(res)
	if err != nil {
		return "", fmt.Errorf("could not marshal json, %v", err)
	}
	return string(b), nil
}

// Config prints the given configuration with weight, amr and bmr
func (r *JSONRenderer) Config(config *model.Config, weight *model.Weight, amr, bmr float64, age int) (string, error) {
	type fullConfig struct {
		Weight     string  `json:"weight"`
		Height     string  `json:"height"`
		Activity   float64 `json:"activity"`
		Birthday   string  `json:"birthday"`
		Age        int     `json:"age"`
		Gender     string  `json:"gender"`
		UnitSystem string  `json:"unitSystem"`
		AMR        float64 `json:"amr"`
		BMR        float64 `json:"bmr"`
	}
	res := fullConfig{
		Weight:     util.WeightUnit(config.UnitSystem, weight.Weight),
		Height:     util.HeightUnit(config.UnitSystem, config.Height),
		Activity:   config.Activity,
		Birthday:   config.Birthday.Format(util.DateFormat),
		Age:        age,
		Gender:     config.Gender,
		UnitSystem: config.UnitSystem,
		AMR:        amr,
		BMR:        bmr,
	}
	b, err := json.Marshal(res)
	if err != nil {
		return "", fmt.Errorf("could not marshal json, %v", err)
	}
	return string(b), nil
}

// Days renders the days in the given timespan
func (r *JSONRenderer) Days(days model.Days, from, to time.Time) (string, error) {
	type daysData struct {
		From time.Time
		To   time.Time
		Days model.Days
	}
	res := daysData{
		From: from,
		To:   to,
		Days: days,
	}
	b, err := json.Marshal(res)
	if err != nil {
		return "", fmt.Errorf("could not marshal json, %v", err)
	}
	return string(b), nil
}

// AddEntry displays a success message after adding an entry
func (r *JSONRenderer) AddEntry(date string, calories int, food string) (string, error) {
	res := success{
		Success: true,
		Message: fmt.Sprintf("Added Entry for %s with %d calories (%s)", date, calories, food),
	}
	b, err := json.Marshal(res)
	if err != nil {
		return "", fmt.Errorf("could not marshal json, %v", err)
	}
	return string(b), nil
}

// ClearEntries displays a success message after clearing the entries for a day
func (r *JSONRenderer) ClearEntries(date string) (string, error) {
	res := success{
		Success: true,
		Message: fmt.Sprintf("Cleared all entries for %s", date),
	}
	b, err := json.Marshal(res)
	if err != nil {
		return "", fmt.Errorf("could not marshal json, %v", err)
	}
	return string(b), nil
}

// ClearEntry displays a success message after clearing the entry at a given position for a day
func (r *JSONRenderer) ClearEntry(date string, entry *model.Entry) (string, error) {
	res := success{
		Success: true,
		Message: fmt.Sprintf("Cleared entry %d %s for %s", entry.Calories, entry.Food, date),
	}
	b, err := json.Marshal(res)
	if err != nil {
		return "", fmt.Errorf("could not marshal json, %v", err)
	}
	return string(b), nil
}

// Export returns the export data as JSON
func (r *JSONRenderer) Export(impex *model.ImpEx) (string, error) {
	b, err := json.MarshalIndent(impex, "", "    ")
	if err != nil {
		return "", fmt.Errorf("could not marshal json, %v", err)
	}
	return string(b), nil
}

// Import displays a success message after importing from a file
func (r *JSONRenderer) Import(fileName string, numEntries, numWeights int) (string, error) {
	res := success{
		Success: true,
		Message: fmt.Sprintf("Imported data from %s with %d entries and %d weights", fileName, numEntries, numWeights),
	}
	b, err := json.Marshal(res)
	if err != nil {
		return "", fmt.Errorf("could not marshal json, %v", err)
	}
	return string(b), nil
}
