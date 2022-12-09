package datasource

import (
	"github.com/zupzup/calories/model"
)

// DataSource is the interface to the data layer
type DataSource interface {
	Setup(connection string) (func() error, error)
	SetConfig(*model.Config) error
	SetConfigFromImport(*model.Config) error
	FetchConfig() (*model.Config, error)
	AddWeight(weight float64) error
	CurrentWeight() (*model.Weight, error)
	FetchWeights() ([]model.Weight, error)
	AddEntry(entryDate string, calories int, food string) error
	FetchEntries(entryDate string) (model.Entries, error)
	FetchAllEntries() (model.Entries, error)
	RemoveEntries(entryDate string) error
	RemoveEntry(entryDate string, id int) error
	Import(data *model.ImpEx) error
	Export() (*model.ImpEx, error)
}
