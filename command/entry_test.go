package command

import (
	"errors"
	"github.com/zupzup/calories/mock"
	"github.com/zupzup/calories/model"
	"testing"
)

func TestExecuteEntriesClearWrongDate(t *testing.T) {
	c := ClearEntriesCommand{
		DataSource: &mock.DataSource{},
		Renderer:   &mock.Renderer{},
		Date:       "bla",
	}
	_, err := c.Execute()
	expected := "wrong format for date: parsing time \"bla\" as \"02.01.2006\": cannot parse \"bla\" as \"02\", please use dd.mm.yyyy"
	if err == nil || err.Error() != expected {
		t.Errorf("Error, actual: %v expected: %v", err.Error(), expected)
		return
	}
}

func TestExecuteEntriesClearZeroEntries(t *testing.T) {
	exps := make(mock.Expectations)
	exps.Add("FetchEntries", nil, model.Entries{})
	c := ClearEntriesCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
		Date:       "01.02.2015",
		Position:   5,
	}
	_, err := c.Execute()
	expected := "could not delete entry at position 5 for 01.02.2015, there are no entries"
	if err == nil || err.Error() != expected {
		t.Errorf("Error, actual: %v expected: %v", err.Error(), expected)
		return
	}
}

func TestExecuteEntriesClearTooFewEntries(t *testing.T) {
	exps := make(mock.Expectations)
	exps.Add("FetchEntries", nil, model.Entries{model.Entry{}})
	c := ClearEntriesCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
		Date:       "01.02.2015",
		Position:   5,
	}
	_, err := c.Execute()
	expected := "could not delete entry at position 5 for 01.02.2015, value needs to be from 1 to 1"
	if err == nil || err.Error() != expected {
		t.Errorf("Error, actual: %v expected: %v", err.Error(), expected)
		return
	}
}

func TestExecuteEntriesClearZeroPosition(t *testing.T) {
	exps := make(mock.Expectations)
	exps.Add("FetchEntries", nil, model.Entries{model.Entry{}})
	c := ClearEntriesCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
		Date:       "01.02.2015",
		Position:   0,
	}
	_, err := c.Execute()
	expected := "could not delete entry at position 0 for 01.02.2015, value needs to be from 1 to 1"
	if err == nil || err.Error() != expected {
		t.Errorf("Error, actual: %v expected: %v", err.Error(), expected)
		return
	}
}

func TestExecuteEntriesClearSuccess(t *testing.T) {
	exps := make(mock.Expectations)
	exps.Add("FetchEntries", nil, model.Entries{model.Entry{}})
	exps.Add("RemoveEntry", nil, nil)
	c := ClearEntriesCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
		Date:       "01.02.2015",
		Position:   1,
		YesMode:    true,
	}
	_, err := c.Execute()
	if err != nil {
		t.Errorf("Error, actual: %v expected: %v", err, nil)
		return
	}
}

func TestExecuteEntriesClearAllSuccess(t *testing.T) {
	exps := make(mock.Expectations)
	exps.Add("RemoveEntries", nil, nil)
	c := ClearEntriesCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
		Date:       "01.02.2015",
		Position:   -1,
		YesMode:    true,
	}
	_, err := c.Execute()
	if err != nil {
		t.Errorf("Error, actual: %v expected: %v", err, nil)
		return
	}
}

func TestExecuteEntryAddWrongUsage(t *testing.T) {
	c := AddEntryCommand{
		DataSource: &mock.DataSource{},
		Renderer:   &mock.Renderer{},
		Mode:       0,
	}
	_, err := c.Execute()
	expected := "usage: calories add [--d=DATE] [--o=FORMAT] CALORIES FOOD"
	if err == nil || err.Error() != expected {
		t.Errorf("Error, actual: %v expected: %v", err.Error(), expected)
		return
	}
}

func TestExecuteEntryAddInvalidCalories(t *testing.T) {
	c := AddEntryCommand{
		DataSource: &mock.DataSource{},
		Renderer:   &mock.Renderer{},
		Mode:       2,
		Calories:   "yay",
	}
	_, err := c.Execute()
	expected := "wrong format for calories: yay needs to be a number (e.g.: 600)"
	if err == nil || err.Error() != expected {
		t.Errorf("Error, actual: %v expected: %v", err.Error(), expected)
		return
	}
}

func TestExecuteEntryAddInvalidDate(t *testing.T) {
	c := AddEntryCommand{
		DataSource: &mock.DataSource{},
		Renderer:   &mock.Renderer{},
		Mode:       2,
		Calories:   "100",
		Date:       "bla",
	}
	_, err := c.Execute()
	expected := "wrong format for date: parsing time \"bla\" as \"02.01.2006\": cannot parse \"bla\" as \"02\", please use dd.mm.yyyy"
	if err == nil || err.Error() != expected {
		t.Errorf("Error, actual: %v expected: %v", err.Error(), expected)
		return
	}
}

func TestExecuteEntryAddEntryFail(t *testing.T) {
	exps := make(mock.Expectations)
	exps.Add("AddEntry", nil, errors.New("someError"))
	c := AddEntryCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
		Mode:       2,
		Calories:   "100",
		Date:       "01.02.2016",
	}
	_, err := c.Execute()
	expected := "someError"
	if err == nil || err.Error() != expected {
		t.Errorf("Error, actual: %v expected: %v", err.Error(), expected)
		return
	}
}

func TestExecuteEntrySuccess(t *testing.T) {
	exps := make(mock.Expectations)
	exps.Add("AddEntry", nil, nil)
	c := AddEntryCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
		Mode:       2,
		Calories:   "100",
		Date:       "01.02.2016",
	}
	_, err := c.Execute()
	if err != nil {
		t.Errorf("Error, actual: %v expected: %v", err, nil)
		return
	}
}
