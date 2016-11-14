package mock

import (
	"github.com/zupzup/calories/model"
	"time"
)

// Renderer is a mocked out renderer
type Renderer struct {
	Expected string
	Err      error
}

// Error Mock
func (r *Renderer) Error(err error) (string, error) {
	return r.Expected, r.Err
}

// WeightHistory Mock
func (r *Renderer) WeightHistory(weights []*model.Weight, config *model.Config) (string, error) {
	return r.Expected, r.Err
}

// AddWeight Mock
func (r *Renderer) AddWeight(weight float64, config *model.Config) (string, error) {
	return r.Expected, r.Err
}

// Config Mock
func (r *Renderer) Config(config *model.Config, weight *model.Weight, amr, bmr float64, age int) (string, error) {
	return r.Expected, r.Err
}

// Days Mock
func (r *Renderer) Days(days model.Days, from, to time.Time) (string, error) {
	return r.Expected, r.Err
}

// AddEntry Mock
func (r *Renderer) AddEntry(date string, calories int, food string) (string, error) {
	return r.Expected, r.Err
}

// ClearEntries Mock
func (r *Renderer) ClearEntries(date string) (string, error) {
	return r.Expected, r.Err
}

// ClearEntry Mock
func (r *Renderer) ClearEntry(date string, entry *model.Entry) (string, error) {
	return r.Expected, r.Err
}

// Import Mock
func (r *Renderer) Import(fileName string, numEntries, numWeights int) (string, error) {
	return r.Expected, r.Err
}
