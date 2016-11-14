package command

import (
	"errors"
	"github.com/zupzup/calories/mock"
	"github.com/zupzup/calories/model"
	"testing"
)

func TestExecuteExportFail(t *testing.T) {
	exps := make(mock.Expectations)
	exps.Add("Export", &model.ImpEx{}, errors.New("err"))
	c := ExportCommand{
		DataSource: &mock.DataSource{Expectations: exps},
	}
	_, err := c.Execute()
	expected := "err"
	if err == nil || err.Error() != expected {
		t.Errorf("Error, actual: %v expected: %v", err.Error(), expected)
		return
	}
}

func TestExecuteExportHappy(t *testing.T) {
	exps := make(mock.Expectations)
	exps.Add("Export", nil, &model.ImpEx{})
	c := ExportCommand{
		DataSource: &mock.DataSource{Expectations: exps},
	}
	res, err := c.Execute()
	if err != nil || res == "" {
		t.Errorf("Error, actual: %v expected: %v", res, "")
		return
	}
}
