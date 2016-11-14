package command

import (
	"fmt"
	"github.com/zupzup/calories/datasource"
	"github.com/zupzup/calories/model"
	"github.com/zupzup/calories/renderer"
	"github.com/zupzup/calories/util"
	"os"
	"time"
)

// ConfigCommand is the command to save and show the configuration
type ConfigCommand struct {
	DataSource datasource.DataSource
	Renderer   renderer.Renderer
	Weight     float64
	Height     float64
	Activity   float64
	Birthday   string
	Gender     string
	UnitSystem string
	YesMode    bool
	Mode       int
}

// Execute shows the current config, if no parameters are given, otherwise it
// parses the given configuration and saves it to the database, asking for
// confirmation first
// The weight from the given config is added to the weight table
func (c *ConfigCommand) Execute() (string, error) {
	if c.Mode < 2 {
		return printConfig(c.DataSource, c.Renderer)
	}
	if c.Weight == -1 || c.Height == -1 || c.Activity == -1 || c.Birthday == "" {
		return "", fmt.Errorf("usage: calories config --w=0.0 --h=0.0 --a=0.0 --b=01.01.1970 --g=male --u=metric")
	}
	parsedBirthday, err := time.Parse(util.DateFormat, c.Birthday)
	if err != nil {
		return "", fmt.Errorf("wrong format for birthday: %v, please use dd.mm.yyyy", err)
	}
	if choice, askErr := checkYesMode(c.YesMode); askErr != nil || !choice {
		return "", askErr
	}
	if err := setConfigAndWeight(c, parsedBirthday); err != nil {
		return "", err
	}
	return printConfig(c.DataSource, c.Renderer)
}

func setConfigAndWeight(c *ConfigCommand, parsedBirthday time.Time) error {
	err := c.DataSource.SetConfig(&model.Config{
		Height:     c.Height,
		Activity:   c.Activity,
		Birthday:   parsedBirthday,
		Gender:     c.Gender,
		UnitSystem: c.UnitSystem,
	})
	if err != nil {
		return fmt.Errorf("could not update config: %v", err)
	}
	err = c.DataSource.AddWeight(c.Weight)
	if err != nil {
		return fmt.Errorf("could not update weight: %v", err)
	}
	return nil
}

func checkYesMode(yesMode bool) (bool, error) {
	if !yesMode {
		return util.AskConfirmation("Do you really want to set this configuration?", os.Stdin)
	}
	return true, nil
}

// printConfig fetches and prints the current config, calculating the age and the metabolic rates
// for the current config
func printConfig(ds datasource.DataSource, r renderer.Renderer) (string, error) {
	config, err := ds.FetchConfig()
	if err != nil {
		return "", fmt.Errorf("could not fetch config: %v", err)
	}
	weight, err := ds.CurrentWeight()
	if err != nil {
		return "", fmt.Errorf("could not fetch current weight: %v", err)
	}
	age := util.CalculateAgeInYears(config.Birthday)
	bmr, amr := util.CalculateHarrisBenedict(float64(age), config.Height, weight.Weight, config.Activity, config.Gender)
	return r.Config(config, weight, amr, bmr, age)
}
