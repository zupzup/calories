package renderer

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/fatih/color"
	"github.com/zupzup/calories/model"
	"github.com/zupzup/calories/util"
)

func TestTerminalError(t *testing.T) {
	r := TerminalRenderer{}
	res, err := r.Error(errors.New("someError"))
	expected := "someError"
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestTerminalWeightHistory(t *testing.T) {
	r := TerminalRenderer{}
	var weights []model.Weight
	now := time.Now()
	weights = append(weights, model.Weight{
		Created: now,
		Weight:  85.0,
	})
	config := model.Config{}
	res, err := r.WeightHistory(weights, &config)
	expected := fmt.Sprintf("Weight over time:\n\t%s: 85.0 kg\n\n", now.Format(util.DateFormat))
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestTerminalAddWeight(t *testing.T) {
	r := TerminalRenderer{}
	res, err := r.AddWeight(85.0, &model.Config{})
	expected := "Set weight: 85.0 kg \n"
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestTerminalConfig(t *testing.T) {
	r := TerminalRenderer{}
	now := time.Now()
	res, err := r.Config(&model.Config{
		Height:     185.0,
		Activity:   1.5,
		Birthday:   now,
		Gender:     "male",
		UnitSystem: util.Metric,
	}, &model.Weight{Weight: 85.0}, 2000.0, 1500.0, 18)
	expected := fmt.Sprintf("Current Config:\n\tWeight: 85.0 kg \n\tHeight: 185.0 cm \n\tActivity: 1.5 \n\tBirthday: %s (18)\n\tGender: male\n\tUnit System: metric\n\tAMR (BMR): 2000 (1500) calories per day\n", now.Format(util.DateFormat))
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestTerminalDaysNoEntries(t *testing.T) {
	r := TerminalRenderer{}
	days := model.Days{}
	now := time.Now()
	res, err := r.Days(days, now, now)
	expected := fmt.Sprintf("Data from %s to %s:\n-----------------------------------\nNo entries have been found.\n", now.Format(util.DateFormat), now.Format(util.DateFormat))
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestTerminalDaysEntries(t *testing.T) {
	r := TerminalRenderer{}
	now := time.Now()
	days := model.Days{}
	entries := model.Entries{}
	entries = append(entries, model.Entry{
		Created:   now,
		EntryDate: now.Format(util.DateFormat),
		Calories:  1000,
		Food:      "Schnitzel",
		BMR:       1500.0,
		AMR:       2000.0,
	})
	days = append(days, &model.Day{
		Used:    1000,
		Date:    now,
		Entries: entries,
	})
	res, err := r.Days(days, now, now)
	expected := fmt.Sprintf("Data from %s to %s:\n-----------------------------------\n%s\n\t1000 Schnitzel\n\t---------------------\n\t%s / 2000 calories\n-----------------------------------\n%s / 2000 calories = %s %s\n", now.Format(util.DateFormat), now.Format(util.DateFormat), now.Format(util.DateFormat), color.GreenString("1000"), color.GreenString("1000"), color.GreenString("1000"), color.GreenString("deficit"))
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestTerminalDaysEntriesSurplus(t *testing.T) {
	r := TerminalRenderer{}
	now := time.Now()
	days := model.Days{}
	entries := model.Entries{}
	entries = append(entries, model.Entry{
		Created:   now,
		EntryDate: now.Format(util.DateFormat),
		Calories:  3000,
		Food:      "Schnitzel",
		BMR:       1500.0,
		AMR:       2000.0,
	})
	days = append(days, &model.Day{
		Used:    3000,
		Date:    now,
		Entries: entries,
	})
	res, err := r.Days(days, now, now)
	expected := fmt.Sprintf("Data from %s to %s:\n-----------------------------------\n%s\n\t3000 Schnitzel\n\t---------------------\n\t%s / 2000 calories\n-----------------------------------\n%s / 2000 calories = %s %s\n", now.Format(util.DateFormat), now.Format(util.DateFormat), now.Format(util.DateFormat), color.RedString("3000"), color.RedString("3000"), color.RedString("1000"), color.RedString("surplus"))
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestTerminalAddEntry(t *testing.T) {
	r := TerminalRenderer{}
	res, err := r.AddEntry("01.01.2017", 1000, "Schnitzel")
	expected := "Added Entry for 01.01.2017 with 1000 calories (Schnitzel)\n"
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestTerminalClearEntries(t *testing.T) {
	r := TerminalRenderer{}
	res, err := r.ClearEntries("01.01.2017")
	expected := "Cleared all entries for 01.01.2017\n"
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestTerminalClearEntry(t *testing.T) {
	r := TerminalRenderer{}
	res, err := r.ClearEntry("01.01.2017", &model.Entry{Calories: 1000.0, Food: "Schnitzel"})
	expected := "Cleared entry 1000 Schnitzel for 01.01.2017\n"
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestTerminalImport(t *testing.T) {
	r := TerminalRenderer{}
	res, err := r.Import("file.csv", 10, 5)
	expected := "Imported data from file.csv with 10 entries and 5 weights\n"
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}
