package command

import (
	"fmt"
	"github.com/zupzup/calories/datasource"
	"github.com/zupzup/calories/renderer"
	"strconv"
)

// WeightCommand is the command to save and show the configuration
type WeightCommand struct {
	DataSource datasource.DataSource
	Renderer   renderer.Renderer
	Weight     string
	Mode       int
}

// Execute shows the weight timeline, if no parameters are given,
// otherwise it sets the given weight
func (c *WeightCommand) Execute() (string, error) {
	config, err := c.DataSource.FetchConfig()
	if err != nil {
		return "", err
	}
	if c.Mode > 0 {
		weight, parseErr := strconv.ParseFloat(c.Weight, 64)
		if parseErr != nil {
			return "", fmt.Errorf("wrong format for weight: %s must be a decimal number", c.Weight)
		}
		err = c.DataSource.AddWeight(weight)
		if err != nil {
			return "", fmt.Errorf("could not save weight: %.0f", weight)
		}
		return c.Renderer.AddWeight(weight, config)
	}

	weights, err := c.DataSource.FetchWeights()
	if err != nil {
		return "", err
	}
	return c.Renderer.WeightHistory(weights, config)
}
