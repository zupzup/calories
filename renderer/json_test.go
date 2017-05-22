package renderer

import (
	"errors"
	"fmt"
	"github.com/zupzup/calories/model"
	"github.com/zupzup/calories/util"
	"testing"
	"time"
)

func TestJSONError(t *testing.T) {
	r := JSONRenderer{}
	res, err := r.Error(errors.New("someError"))
	expected := fmt.Sprint("{\"success\":false,\"message\":\"someError\"}")
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestJSONWeightHistory(t *testing.T) {
	r := JSONRenderer{}
	var weights []model.Weight
	now := time.Now()
	weights = append(weights, model.Weight{
		Created: now,
		Weight:  85.0,
	})
	config := model.Config{}
	res, err := r.WeightHistory(weights, &config)
	expected := fmt.Sprintf("[{\"created\":\"%s\",\"weight\":85,\"formatted\":\"85.0 kg\"}]", now.Format(time.RFC3339Nano))
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestJSONConfig(t *testing.T) {
	r := JSONRenderer{}
	now := time.Now()
	res, err := r.Config(&model.Config{
		Height:     185.0,
		Activity:   1.5,
		Birthday:   now,
		Gender:     "male",
		UnitSystem: util.Metric,
	}, &model.Weight{Weight: 85.0}, 2000.0, 1500.0, 18)
	expected := fmt.Sprintf("{\"weight\":\"85.0 kg\",\"height\":\"185.0 cm\",\"activity\":1.5,\"birthday\":\"%s\",\"age\":18,\"gender\":\"male\",\"unitSystem\":\"metric\",\"amr\":2000,\"bmr\":1500}", now.Format(util.DateFormat))
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestJSONAddWeight(t *testing.T) {
	r := JSONRenderer{}
	res, err := r.AddWeight(85.0, &model.Config{})
	expectedString := "Added weight: 85.0 kg"
	expected := fmt.Sprintf("{\"success\":true,\"message\":\"%s\"}", expectedString)
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestJSONDaysNoEntries(t *testing.T) {
	r := JSONRenderer{}
	days := model.Days{}
	now := time.Now()
	res, err := r.Days(days, now, now)
	expected := fmt.Sprintf("{\"From\":\"%s\",\"To\":\"%s\",\"Days\":[]}", now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano))
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestJSONDaysEntries(t *testing.T) {
	r := JSONRenderer{}
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
	expected := fmt.Sprintf("{\"From\":\"%s\",\"To\":\"%s\",\"Days\":[{\"entries\":[{\"id\":0,\"created\":\"%s\",\"entryDate\":\"%s\",\"calories\":1000,\"food\":\"Schnitzel\",\"bmr\":1500,\"amr\":2000}],\"used\":1000,\"date\":\"%s\"}]}", now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano), now.Format(util.DateFormat), now.Format(time.RFC3339Nano))
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestJSONAddEntry(t *testing.T) {
	r := JSONRenderer{}
	res, err := r.AddEntry("01.01.2017", 1000, "Schnitzel")
	expectedString := "Added Entry for 01.01.2017 with 1000 calories (Schnitzel)"
	expected := fmt.Sprintf("{\"success\":true,\"message\":\"%s\"}", expectedString)
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestJSONClearEntries(t *testing.T) {
	r := JSONRenderer{}
	res, err := r.ClearEntries("01.01.2017")
	expectedString := "Cleared all entries for 01.01.2017"
	expected := fmt.Sprintf("{\"success\":true,\"message\":\"%s\"}", expectedString)
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestJSONClearEntry(t *testing.T) {
	r := JSONRenderer{}
	res, err := r.ClearEntry("01.01.2017", &model.Entry{Calories: 1000.0, Food: "Schnitzel"})
	expectedString := "Cleared entry 1000 Schnitzel for 01.01.2017"
	expected := fmt.Sprintf("{\"success\":true,\"message\":\"%s\"}", expectedString)
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestJSONImport(t *testing.T) {
	r := JSONRenderer{}
	res, err := r.Import("file.csv", 10, 5)
	expectedString := "Imported data from file.csv with 10 entries and 5 weights"
	expected := fmt.Sprintf("{\"success\":true,\"message\":\"%s\"}", expectedString)
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}
