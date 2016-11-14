package util

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"time"
)

func TestAskConfirmationPassYes(t *testing.T) {
	yReader := strings.NewReader("yes\n")
	res, err := AskConfirmation("hi", yReader)
	expected := true
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestAskConfirmationPassNo(t *testing.T) {
	yReader := strings.NewReader("no\n")
	res, err := AskConfirmation("hi", yReader)
	expected := false
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestAskConfirmationFailEOF(t *testing.T) {
	yReader := strings.NewReader("yes")
	res, err := AskConfirmation("hi", yReader)
	expected := false
	if res != expected || err == nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

var testsSort = []struct {
	description string
	in          time.Time
	result      int
}{
	{
		"basic",
		time.Now().AddDate(-10, 0, 0),
		10,
	},
	{
		"same",
		time.Now(),
		0,
	},
	{
		"month greater",
		time.Now().AddDate(-10, 1, 0),
		9,
	},
	{
		"day greater",
		time.Now().AddDate(-10, 0, 1),
		9,
	},
	{
		"month smaller",
		time.Now().AddDate(-10, -1, 0),
		10,
	},
	{
		"day smaller",
		time.Now().AddDate(-10, 0, -1),
		10,
	},
}

func TestCalculateAgeInYears(t *testing.T) {
	for _, tc := range testsSort {
		t.Run(fmt.Sprintf("Test: %s", tc.description), func(t *testing.T) {
			res := CalculateAgeInYears(tc.in)
			if res != tc.result {
				t.Errorf("Error, actual: %v expected: %v", res, tc.result)
				return
			}
		})
	}
}

func TestGetBeginningOfWeek(t *testing.T) {
	date, _ := time.Parse(DateFormat, "15.11.2016")
	expected, _ := time.Parse(DateFormat, "14.11.2016")
	res := GetBeginningOfWeek(date)
	if !res.Equal(expected) {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestGetBeginningOfWeekSunday(t *testing.T) {
	date, _ := time.Parse(DateFormat, "06.11.2016")
	expected, _ := time.Parse(DateFormat, "31.10.2016")
	res := GetBeginningOfWeek(date)
	if !res.Equal(expected) {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestCalculateHarrisBenedict(t *testing.T) {
	bmr, amr := CalculateHarrisBenedict(0, 0, 0, 1.3, "male")
	fbmr, famr := CalculateHarrisBenedict(0, 0, 0, 1, "female")
	if bmr != 88.362 || amr != 114.8706 || fbmr != 447.593 || famr != 447.593 {
		t.Errorf("Error, actual: %v %v %v %v expected: %v %v %v %v", bmr, amr, fbmr, famr, 66.5, 86.45, 655.1, 655.1)
		return
	}
}

var testsWeightUnit = []struct {
	description string
	unit        string
	value       float64
	result      string
}{
	{
		"imp",
		Imperial,
		10.0,
		"22.0 pounds",
	},
	{
		"metric",
		Metric,
		10.0,
		"10.0 kg",
	},
}

func TestWeightUnit(t *testing.T) {
	for _, tc := range testsWeightUnit {
		t.Run(fmt.Sprintf("Test: %s", tc.description), func(t *testing.T) {
			res := WeightUnit(tc.unit, tc.value)
			if res != tc.result {
				t.Errorf("Error, actual: %v expected: %v", res, tc.result)
				return
			}
		})
	}
}

var testsHeightUnit = []struct {
	description string
	unit        string
	value       float64
	result      string
}{
	{
		"imp",
		Imperial,
		185.0,
		"72.8 inches",
	},
	{
		"metric",
		Metric,
		185.0,
		"185.0 cm",
	},
}

func TestHeightUnit(t *testing.T) {
	for _, tc := range testsHeightUnit {
		t.Run(fmt.Sprintf("Test: %s", tc.description), func(t *testing.T) {
			res := HeightUnit(tc.unit, tc.value)
			if res != tc.result {
				t.Errorf("Error, actual: %v expected: %v", res, tc.result)
				return
			}
		})
	}
}

func TestConversionsToMetric(t *testing.T) {
	kg := math.Ceil(ToKg(22.0))
	cm := math.Ceil(ToCm(72.8))
	if kg != 10.0 || cm != 185.0 {
		t.Errorf("Error, actual: %v %v expected: %v %v", kg, cm, 10.0, 185.0)
		return
	}
}
