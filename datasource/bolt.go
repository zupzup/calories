package datasource

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/zupzup/calories/model"
	"github.com/zupzup/calories/util"
	"time"
)

// BoltDataSource is an implementation of the DataSource interface for boltdb
type BoltDataSource struct {
	DB *storm.DB
}

// Setup creates the file and the table structure
func (ds *BoltDataSource) Setup(connection string) (func() error, error) {
	db, err := storm.Open(connection)
	if err != nil {
		return nil, fmt.Errorf("error while connecting to database at %s, %v", connection, err)
	}
	ds.DB = db
	return db.Close, nil
}

// SetConfig overrides the current config with the given values
// by deleting the old config and adding a new one
func (ds *BoltDataSource) SetConfig(c *model.Config) error {
	ds.DB.Drop(&model.Config{})
	if c.UnitSystem != "metric" && c.UnitSystem != "imperial" {
		return fmt.Errorf("unit system needs to be either metric or imperial: %s", c.UnitSystem)
	}
	height := c.Height
	if c.UnitSystem == util.Imperial {
		height = util.ToCm(height)
	}
	config := model.Config{
		Height:     height,
		Activity:   c.Activity,
		Birthday:   c.Birthday,
		Gender:     c.Gender,
		UnitSystem: c.UnitSystem,
	}
	err := ds.DB.Save(&config)
	return err
}

// FetchConfig fetches and returns the current config
func (ds *BoltDataSource) FetchConfig() (config *model.Config, err error) {
	var configs []model.Config
	err = ds.DB.All(&configs, storm.Limit(1), storm.Reverse())
	if err != nil || len(configs) == 0 {
		return nil, fmt.Errorf("could not retrieve config: %v", err)
	}
	return &configs[0], nil
}

// AddWeight adds the given weight for todays date
func (ds *BoltDataSource) AddWeight(weight float64) error {
	config, err := ds.FetchConfig()
	if err != nil {
		return err
	}
	if config.UnitSystem == util.Imperial {
		weight = util.ToKg(weight)
	}
	weightObj := model.Weight{
		Created: time.Now(),
		Weight:  weight,
	}
	err = ds.DB.Save(&weightObj)
	return err
}

// CurrentWeight fetches and returns the current weight, which is the last entry in the table
func (ds *BoltDataSource) CurrentWeight() (*model.Weight, error) {
	var weights []model.Weight
	err := ds.DB.All(&weights, storm.Limit(1), storm.Reverse())
	if err != nil || len(weights) == 0 {
		return nil, fmt.Errorf("could not fetch current weight: %v", err)
	}
	return &weights[0], nil
}

// FetchWeights fetches all weight entries
func (ds *BoltDataSource) FetchWeights() ([]model.Weight, error) {
	var weights []model.Weight
	err := ds.DB.All(&weights)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve weight history: %v", err)
	}
	return weights, nil
}

// AddEntry fetches the current config and weight to calculate the metabolic rate and adds the data
// into the entry table
func (ds *BoltDataSource) AddEntry(entryDate string, calories int, food string) error {
	weight, err := ds.CurrentWeight()
	if err != nil {
		return err
	}
	config, err := ds.FetchConfig()
	if err != nil {
		return err
	}
	age := float64(util.CalculateAgeInYears(config.Birthday))
	bmr, amr := util.CalculateHarrisBenedict(age, config.Height, weight.Weight, config.Activity, config.Gender)
	entry := model.Entry{
		Created:   time.Now(),
		EntryDate: entryDate,
		Calories:  calories,
		Food:      food,
		AMR:       amr,
		BMR:       bmr,
	}
	err = ds.DB.Save(&entry)
	if err != nil {
		return fmt.Errorf("could not add entry: %v", err)
	}
	return nil
}

// FetchEntries fetches and returns all entries for a given date
func (ds *BoltDataSource) FetchEntries(entryDate string) (model.Entries, error) {
	var entries []model.Entry
	err := ds.DB.Find("EntryDate", entryDate, &entries)
	if err != nil {
		if err == storm.ErrNotFound {
			return entries, nil
		}
		return nil, fmt.Errorf("could not fetch entries for the given date: %s, %v", entryDate, err)
	}
	return entries, nil
}

// FetchAllEntries fetches and returns all entries
func (ds *BoltDataSource) FetchAllEntries() (model.Entries, error) {
	var entries []model.Entry
	err := ds.DB.All(&entries)
	if err != nil {
		if err == storm.ErrNotFound {
			return entries, nil
		}
		return nil, fmt.Errorf("could not fetch all entries, %v", err)
	}
	return entries, nil
}

// RemoveEntries removes all entries for a given day from the database
func (ds *BoltDataSource) RemoveEntries(entryDate string) error {
	query := ds.DB.Select(q.Eq("EntryDate", entryDate))
	err := query.Delete(new(model.Entry))
	if err != nil {
		return fmt.Errorf("could not delete entries for %s", entryDate)
	}
	return nil
}

// RemoveEntry removes the entry with the given id for a given day from the database
func (ds *BoltDataSource) RemoveEntry(entryDate string, id int) error {
	query := ds.DB.Select(q.And(q.Eq("EntryDate", entryDate), q.Eq("ID", id)))
	err := query.Delete(new(model.Entry))
	if err != nil {
		return fmt.Errorf("could not delete entry with id %d on day %s", id, entryDate)
	}
	return nil
}

// Import imports the given data to the database, overwriting the previous
// data
func (ds *BoltDataSource) Import(data *model.ImpEx) error {
	err := ds.SetConfig(data.Config)
	if err != nil {
		return fmt.Errorf("could not replace config, %v", err)
	}
	err = ds.DB.Drop(&model.Weight{})
	if err != nil {
		return fmt.Errorf("could not remove weights, %v", err)
	}
	for _, weight := range data.Weights {
		err = ds.DB.Save(&weight)
		if err != nil {
			return fmt.Errorf("could not insert/update weight with id %d", weight.ID)
		}
	}
	err = ds.DB.Drop(&model.Entry{})
	if err != nil {
		return fmt.Errorf("could not remove entries, %v", err)
	}
	for _, entry := range data.Entries {
		err = ds.DB.Save(&entry)
		if err != nil {
			return fmt.Errorf("could not insert/update entry with id %d", entry.ID)
		}
	}
	return nil
}

// Export creates a JSON representation of the database
func (ds *BoltDataSource) Export() (*model.ImpEx, error) {
	entries, err := ds.FetchAllEntries()
	if err != nil {
		return nil, err
	}
	weights, err := ds.FetchWeights()
	if err != nil {
		return nil, err
	}
	config, err := ds.FetchConfig()
	if err != nil {
		return nil, err
	}
	impex := &model.ImpEx{
		Config:  config,
		Entries: entries,
		Weights: weights,
	}
	return impex, nil
}
