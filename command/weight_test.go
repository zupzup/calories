package command

import (
	"errors"
	"github.com/zupzup/calories/mock"
	"github.com/zupzup/calories/model"
	"testing"
)

func TestExecuteWeightNoAdd(t *testing.T) {
	var weights []*model.Weight
	exps := make(mock.Expectations)
	exps.Add("FetchConfig", nil, &model.Config{})
	exps.Add("FetchWeights", nil, weights)
	c := WeightCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
	}
	res, err := c.Execute()
	expected := ""
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestExecuteWeightBrokenWeights(t *testing.T) {
	exps := make(mock.Expectations)
	exps.Add("FetchConfig", nil, &model.Config{})
	exps.Add("FetchWeights", []*model.Weight{}, errors.New("err"))
	c := WeightCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
	}
	_, err := c.Execute()
	if err == nil {
		t.Errorf("Error, actual: %v expected: %v", err, "err")
		return
	}
}

func TestExecuteWeightNoConfig(t *testing.T) {
	exps := make(mock.Expectations)
	exps.Add("FetchConfig", &model.Config{}, errors.New("err"))
	c := WeightCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
	}
	_, err := c.Execute()
	if err == nil {
		t.Errorf("Error, actual: %v expected: %v", err, "err")
		return
	}
}

func TestExecuteWeightAdd(t *testing.T) {
	var weights []*model.Weight
	exps := make(mock.Expectations)
	exps.Add("FetchConfig", nil, &model.Config{})
	exps.Add("AddWeight", nil, nil)
	exps.Add("FetchWeights", nil, weights)
	c := WeightCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
		Mode:       1,
		Weight:     "85",
	}
	res, err := c.Execute()
	expected := ""
	if res != expected || err != nil {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}

func TestExecuteWeightAddNoNumber(t *testing.T) {
	var weights []*model.Weight
	exps := make(mock.Expectations)
	exps.Add("FetchConfig", nil, &model.Config{})
	exps.Add("AddWeight", nil, nil)
	exps.Add("FetchWeights", nil, weights)
	c := WeightCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
		Mode:       1,
		Weight:     "hello",
	}
	_, err := c.Execute()
	expected := "wrong format for weight: hello must be a decimal number"
	if err == nil || err.Error() != expected {
		t.Errorf("Error, actual: %v expected: %v", err.Error(), expected)
		return
	}
}

func TestExecuteWeightAddBroken(t *testing.T) {
	var weights []*model.Weight
	exps := make(mock.Expectations)
	exps.Add("FetchConfig", nil, &model.Config{})
	exps.Add("AddWeight", nil, errors.New("err"))
	exps.Add("FetchWeights", nil, weights)
	c := WeightCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
		Mode:       1,
		Weight:     "85",
	}
	_, err := c.Execute()
	expected := "could not save weight: 85"
	if err == nil || err.Error() != expected {
		t.Errorf("Error, actual: %v expected: %v", err.Error(), expected)
		return
	}
}
