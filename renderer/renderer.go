package renderer

import (
	"github.com/zupzup/calories/model"
	"time"
)

// Renderer is the interface for rendering any output
type Renderer interface {
	Error(err error) (string, error)
	WeightHistory(weights []model.Weight, config *model.Config) (string, error)
	AddWeight(weight float64, config *model.Config) (string, error)
	Config(config *model.Config, weight *model.Weight, amr, bmr float64, age int) (string, error)
	Days(days model.Days, from, to time.Time) (string, error)
	AddEntry(date string, calories int, food string) (string, error)
	ClearEntries(date string) (string, error)
	ClearEntry(date string, entry *model.Entry) (string, error)
	Import(fileName string, numEntries, numWeights int) (string, error)
}
