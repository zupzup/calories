package command

import (
	"github.com/zupzup/calories/datasource"
	"github.com/zupzup/calories/renderer"
)

// ExportCommand is the command to export the database
type ExportCommand struct {
	DataSource datasource.DataSource
}

// Execute fetches the export data and writes it to the given file as JSON
func (c *ExportCommand) Execute() (string, error) {
	impex, err := c.DataSource.Export()
	if err != nil {
		return "", err
	}
	jsonRenderer := &renderer.JSONRenderer{}
	return jsonRenderer.Export(impex)
}
