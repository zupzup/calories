package command

import (
	"encoding/json"
	"fmt"
	"github.com/zupzup/calories/datasource"
	"github.com/zupzup/calories/model"
	"github.com/zupzup/calories/renderer"
	"io/ioutil"
)

// ImportCommand is the command to export the database
type ImportCommand struct {
	DataSource datasource.DataSource
	Renderer   renderer.Renderer
	File       string
}

// Execute parses and imports the data from the given file
func (c *ImportCommand) Execute() (string, error) {
	if c.File == "" {
		return "", fmt.Errorf("no import file provided")
	}
	b, err := ioutil.ReadFile(c.File)
	if err != nil {
		return "", fmt.Errorf("error reading file %s, %v", c.File, err)
	}
	impex := model.ImpEx{}
	err = json.Unmarshal(b, &impex)
	if err != nil {
		return "", fmt.Errorf("error parsing json, %v", err)
	}
	err = c.DataSource.Import(&impex)
	if err != nil {
		return "", err
	}
	return c.Renderer.Import(c.File, len(impex.Entries), len(impex.Weights))
}
