package command

import (
	"errors"
	"github.com/zupzup/calories/mock"
	"testing"
)

func TestExecuteImportNoFile(t *testing.T) {
	c := ImportCommand{DataSource: &mock.DataSource{}, Renderer: &mock.Renderer{}, File: ""}
	_, err := c.Execute()
	expected := "no import file provided"
	if err == nil || err.Error() != expected {
		t.Errorf("Error, actual: %v expected: %v", err.Error(), expected)
		return
	}
}

func TestExecuteImportInvalidFile(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	c := ImportCommand{DataSource: &mock.DataSource{}, Renderer: &mock.Renderer{}, File: "test.txt"}
	_, err := c.Execute()
	expected := "error reading file test.txt, open test.txt: no such file or directory"
	if err == nil || err.Error() != expected {
		t.Errorf("Error, actual: %v expected: %v", err.Error(), expected)
		return
	}
}

func TestExecuteImportInvalidJSONFile(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	c := ImportCommand{DataSource: &mock.DataSource{}, Renderer: &mock.Renderer{}, File: "./testdata/import.txt"}
	_, err := c.Execute()
	expected := "error parsing json, unexpected end of JSON input"
	if err == nil || err.Error() != expected {
		t.Errorf("Error, actual: %v expected: %v", err.Error(), expected)
		return
	}
}

func TestExecuteImportFail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	exps := make(mock.Expectations)
	exps.Add("Import", nil, errors.New("err"))
	c := ImportCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
		File:       "./testdata/import.json",
	}
	_, err := c.Execute()
	expected := "err"
	if err == nil || err.Error() != expected {
		t.Errorf("Error, actual: %v expected: %v", err.Error(), expected)
		return
	}
}

func TestExecuteImportHappy(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	exps := make(mock.Expectations)
	exps.Add("Import", nil, nil)
	c := ImportCommand{
		DataSource: &mock.DataSource{Expectations: exps},
		Renderer:   &mock.Renderer{},
		File:       "./testdata/import.json",
	}
	res, err := c.Execute()
	expected := ""
	if err != nil || res != expected {
		t.Errorf("Error, actual: %v expected: %v", res, expected)
		return
	}
}
