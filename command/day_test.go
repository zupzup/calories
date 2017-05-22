package command

import (
	"errors"
	"github.com/zupzup/calories/mock"
	"github.com/zupzup/calories/model"
	"testing"
)

func TestExecuteDayWeekSuccessEmpty(t *testing.T) {
	exps := make(mock.Expectations)
	for i := 0; i <= 31; i++ {
		exps.Add("FetchEntries", nil, model.Entries{model.Entry{}})
	}
	c := DayCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
		Week:       true,
	}
	_, err := c.Execute()
	if err != nil {
		t.Errorf("Error, actual: %v expected: %v", err, nil)
		return
	}
}

func TestExecuteDayWeekFetchFail(t *testing.T) {
	exps := make(mock.Expectations)
	for i := 0; i <= 31; i++ {
		exps.Add("FetchEntries", model.Entries{}, errors.New("someError"))
	}
	c := DayCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
		Week:       true,
	}
	_, err := c.Execute()
	expected := "someError"
	if err == nil || err.Error() != expected {
		t.Errorf("Error, actual: %v expected: %v", err.Error(), expected)
		return
	}
}

func TestExecuteDayWeekSuccessEntries(t *testing.T) {
	exps := make(mock.Expectations)
	for i := 0; i <= 31; i++ {
		exps.Add("FetchEntries", nil, model.Entries{model.Entry{}})
	}
	c := DayCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
		Week:       true,
	}
	_, err := c.Execute()
	if err != nil {
		t.Errorf("Error, actual: %v expected: %v", err, nil)
		return
	}
}

func TestExecuteDayMonthSuccess(t *testing.T) {
	exps := make(mock.Expectations)
	for i := 0; i <= 31; i++ {
		exps.Add("FetchEntries", nil, model.Entries{model.Entry{}})
	}
	c := DayCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
		Week:       false,
		Month:      true,
	}
	_, err := c.Execute()
	if err != nil {
		t.Errorf("Error, actual: %v expected: %v", err, nil)
		return
	}
}

func TestExecuteDayDateSuccess(t *testing.T) {
	exps := make(mock.Expectations)
	exps.Add("FetchEntries", nil, model.Entries{model.Entry{}})
	c := DayCommand{
		DataSource:  &mock.DataSource{Expectations: exps},
		Renderer:    &mock.Renderer{},
		Week:        false,
		Month:       false,
		History:     0,
		DefaultDate: "02.01.2016",
	}
	_, err := c.Execute()
	if err != nil {
		t.Errorf("Error, actual: %v expected: %v", err, nil)
		return
	}
}

func TestExecuteDayFalseDate(t *testing.T) {
	exps := make(mock.Expectations)
	exps.Add("FetchEntries", nil, model.Entries{model.Entry{}})
	c := DayCommand{
		DataSource:  &mock.DataSource{Expectations: exps},
		Renderer:    &mock.Renderer{},
		Week:        false,
		Month:       false,
		History:     0,
		DefaultDate: "bla",
	}
	_, err := c.Execute()
	expected := "wrong format for date: parsing time \"bla\" as \"02.01.2006\": cannot parse \"bla\" as \"02\", please use dd.mm.yyyy"
	if err == nil || err.Error() != expected {
		t.Errorf("Error, actual: %v expected: %v", err.Error(), expected)
		return
	}
}

func TestExecuteDayNoDateSuccess(t *testing.T) {
	exps := make(mock.Expectations)
	exps.Add("FetchEntries", nil, model.Entries{model.Entry{}})
	c := DayCommand{
		DataSource:  &mock.DataSource{Expectations: exps},
		Renderer:    &mock.Renderer{},
		Week:        false,
		Month:       false,
		History:     0,
		DefaultDate: "",
	}
	_, err := c.Execute()
	if err != nil {
		t.Errorf("Error, actual: %v expected: %v", err, nil)
		return
	}
}

func TestExecuteDayHistorySuccess(t *testing.T) {
	exps := make(mock.Expectations)
	for i := 0; i <= 6; i++ {
		exps.Add("FetchEntries", nil, model.Entries{model.Entry{}})
	}
	c := DayCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
		Week:       false,
		Month:      false,
		History:    5,
	}
	_, err := c.Execute()
	if err != nil {
		t.Errorf("Error, actual: %v expected: %v", err, nil)
		return
	}
}

func TestExecuteDayHistoryMinusSuccess(t *testing.T) {
	exps := make(mock.Expectations)
	for i := 0; i <= 6; i++ {
		exps.Add("FetchEntries", nil, model.Entries{model.Entry{}})
	}
	c := DayCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
		Week:       false,
		Month:      false,
		History:    -5,
	}
	_, err := c.Execute()
	if err != nil {
		t.Errorf("Error, actual: %v expected: %v", err, nil)
		return
	}
}
