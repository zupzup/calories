package util

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

// DateFormat is the date format we are using
const DateFormat = "02.01.2006"

// Imperial depicts the identifier for the imperial unit system
const Imperial = "imperial"

// Metric depicts the identifier for the metric unit system
const Metric = "metric"

// AskConfirmation asks the user for confirmation on a given question
// and returns the user's answer
func AskConfirmation(s string, r io.Reader) (bool, error) {
	prompt := bufio.NewReader(r)
	for {
		fmt.Printf("%s [y/n]: ", s)
		response, err := prompt.ReadString('\n')
		if err != nil {
			return false, err
		}
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return true, nil
		} else if response == "n" || response == "no" {
			return false, nil
		}
	}
}

// CalculateHarrisBenedict calculates Harris-Benedict (https://en.wikipedia.org/wiki/Harris%E2%80%93Benedict_equation)
// for calculating the basic metabolic rate
func CalculateHarrisBenedict(age, height, weight, activity float64, gender string) (float64, float64) {
	baseLine := 88.362
	weightMultiplier := 13.397
	heightMultiplier := 4.799
	ageMultiplier := 5.677
	if gender == "female" {
		baseLine = 447.593
		weightMultiplier = 9.247
		heightMultiplier = 3.098
		ageMultiplier = 4.330
	}
	basicMetabolicRate := (baseLine + weightMultiplier*weight + heightMultiplier*height - ageMultiplier*age)
	return basicMetabolicRate, basicMetabolicRate * activity
}

// CalculateAgeInYears calculates the age in years given a date by comparing the
// year, month and day of now and the given date
func CalculateAgeInYears(birthday time.Time) int {
	now := time.Now()
	yearsDiff := now.Year() - birthday.Year()
	if now.Month() < birthday.Month() {
		return yearsDiff - 1
	}
	if now.Month() > birthday.Month() || now.Day() >= birthday.Day() {
		return yearsDiff
	}
	return yearsDiff - 1
}

// GetBeginningOfWeek calculates the first day of the week (Monday) given a date
func GetBeginningOfWeek(date time.Time) time.Time {
	mondayDiff := -int(date.Weekday()) + 1
	if int(date.Weekday()) == 0 {
		mondayDiff = -6
	}
	return date.AddDate(0, 0, mondayDiff)
}

// WeightUnit returns the weight unit for the given unit system
func WeightUnit(unitSystem string, value float64) string {
	unit := "kg"
	if unitSystem == Imperial {
		value = ToPounds(value)
		unit = "pounds"
	}
	return fmt.Sprintf("%.1f %s", value, unit)
}

// HeightUnit returns the height unit for the given unit system
func HeightUnit(unitSystem string, value float64) string {
	unit := "cm"
	if unitSystem == Imperial {
		value = ToInches(value)
		unit = "inches"
	}
	return fmt.Sprintf("%.1f %s", value, unit)
}

// ToPounds converts kg to pounds
func ToPounds(weight float64) float64 {
	return weight * 2.20462
}

// ToInches converts cm to inches
func ToInches(height float64) float64 {
	return height * 0.393701
}

// ToKg converts pounds to kg
func ToKg(weight float64) float64 {
	return weight / 2.20462
}

// ToCm converts inches to cm
func ToCm(height float64) float64 {
	return height / 0.393701
}
