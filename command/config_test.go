package command

import (
	"errors"
	"github.com/zupzup/calories/mock"
	"github.com/zupzup/calories/model"
	"testing"
	"time"
)

var dummyConfig = model.Config{
	Height:     101.5,
	Birthday:   time.Now(),
	Activity:   1.2,
	Gender:     "male",
	UnitSystem: "metric",
}

var dummyWeight = model.Weight{
	Weight: 85.9,
}

func TestExecuteConfigOutputModeFetchConfigFail(t *testing.T) {
	exps := make(mock.Expectations)
	exps.Add("FetchConfig", &model.Config{}, errors.New("someError"))
	c := ConfigCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
		Mode:       0,
	}
	_, err := c.Execute()
	expected := "could not fetch config: someError"
	if err == nil || err.Error() != expected {
		t.Errorf("Error, actual: %v expected: %v", err.Error(), expected)
		return
	}
}

func TestExecuteConfigOutputModeFetchWeightFail(t *testing.T) {
	exps := make(mock.Expectations)
	exps.Add("FetchConfig", nil, &model.Config{})
	exps.Add("CurrentWeight", &model.Weight{}, errors.New("someError"))
	c := ConfigCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
		Mode:       0,
	}
	_, err := c.Execute()
	expected := "could not fetch current weight: someError"
	if err == nil || err.Error() != expected {
		t.Errorf("Error, actual: %v expected: %v", err.Error(), expected)
		return
	}
}

func TestExecuteConfigOutputModeSuccess(t *testing.T) {
	exps := make(mock.Expectations)
	exps.Add("FetchConfig", nil, &dummyConfig)
	exps.Add("CurrentWeight", nil, &dummyWeight)
	c := ConfigCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
		Mode:       0,
	}
	_, err := c.Execute()
	if err != nil {
		t.Errorf("Error, actual: %v expected: %v", err, nil)
		return
	}
}

func TestExecuteConfigSetModeInvalid(t *testing.T) {
	c := ConfigCommand{
		DataSource: &mock.DataSource{},
		Renderer:   &mock.Renderer{},
		Mode:       2,
	}
	_, err := c.Execute()
	expected := "usage: calories config --w=0.0 --h=0.0 --a=0.0 --b=01.01.1970 --g=male --u=metric"
	if err == nil || err.Error() != expected {
		t.Errorf("Error, actual: %v expected: %v", err.Error(), expected)
		return
	}
}

func TestExecuteConfigSetModeInvalidBirthday(t *testing.T) {
	c := ConfigCommand{
		DataSource: &mock.DataSource{},
		Renderer:   &mock.Renderer{},
		Mode:       2,
		Weight:     85.0,
		Height:     185.9,
		Activity:   1.3,
		Birthday:   "bla",
	}
	_, err := c.Execute()
	expected := "wrong format for birthday: parsing time \"bla\" as \"02.01.2006\": cannot parse \"bla\" as \"02\", please use dd.mm.yyyy"
	if err == nil || err.Error() != expected {
		t.Errorf("Error, actual: %v expected: %v", err.Error(), expected)
		return
	}
}

func TestExecuteConfigSetModeSuccess(t *testing.T) {
	exps := make(mock.Expectations)
	exps.Add("FetchConfig", nil, &dummyConfig)
	exps.Add("CurrentWeight", nil, &dummyWeight)
	exps.Add("SetConfig", nil, nil)
	exps.Add("AddWeight", nil, nil)
	c := ConfigCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
		Mode:       2,
		Weight:     85.0,
		Height:     185.9,
		Activity:   1.3,
		Birthday:   "08.08.1985",
		YesMode:    true,
	}
	_, err := c.Execute()
	if err != nil {
		t.Errorf("Error, actual: %v expected: %v", err, nil)
		return
	}
}
