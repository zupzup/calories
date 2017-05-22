package mock

import (
	"github.com/zupzup/calories/model"
)

// DataSource is a mocked out datasource
type DataSource struct {
	Expectations Expectations
}

// Setup Mock
func (d *DataSource) Setup(connection string) (func() error, error) {
	v, err := d.Expectations.Return("Setup")
	return v.(func() error), err
}

// SetConfig Mock
func (d *DataSource) SetConfig(*model.Config) error {
	_, err := d.Expectations.Return("SetConfig")
	return err
}

// FetchConfig Mock
func (d *DataSource) FetchConfig() (*model.Config, error) {
	v, err := d.Expectations.Return("FetchConfig")
	return v.(*model.Config), err
}

// AddWeight Mock
func (d *DataSource) AddWeight(weight float64) error {
	_, err := d.Expectations.Return("AddWeight")
	return err
}

// CurrentWeight Mock
func (d *DataSource) CurrentWeight() (*model.Weight, error) {
	v, err := d.Expectations.Return("CurrentWeight")
	return v.(*model.Weight), err
}

// FetchWeights Mock
func (d *DataSource) FetchWeights() ([]model.Weight, error) {
	v, err := d.Expectations.Return("FetchWeights")
	return v.([]model.Weight), err
}

// AddEntry Mock
func (d *DataSource) AddEntry(entryDate string, calories int, food string) error {
	_, err := d.Expectations.Return("AddEntry")
	return err
}

// FetchEntries Mock
func (d *DataSource) FetchEntries(entryDate string) (model.Entries, error) {
	v, err := d.Expectations.Return("FetchEntries")
	return v.(model.Entries), err
}

// FetchAllEntries Mock
func (d *DataSource) FetchAllEntries() (model.Entries, error) {
	v, err := d.Expectations.Return("FetchAllEntries")
	return v.(model.Entries), err
}

// RemoveEntries Mock
func (d *DataSource) RemoveEntries(entryDate string) error {
	_, err := d.Expectations.Return("RemoveEntries")
	return err
}

// RemoveEntry Mock
func (d *DataSource) RemoveEntry(entryDate string, id int) error {
	_, err := d.Expectations.Return("RemoveEntry")
	return err
}

// Import Mock
func (d *DataSource) Import(data *model.ImpEx) error {
	_, err := d.Expectations.Return("Import")
	return err
}

// Export Mock
func (d *DataSource) Export() (*model.ImpEx, error) {
	v, err := d.Expectations.Return("Export")
	return v.(*model.ImpEx), err
}
