package datasource

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // sqlite driver
	"github.com/zupzup/calories/model"
	"github.com/zupzup/calories/util"
	"strings"
	"time"
)

// SQLiteDataSource is an implementation of the DataSource interface for sqlite
type SQLiteDataSource struct {
	DB *sql.DB
}

// Setup creates the file and the table structure
func (ds *SQLiteDataSource) Setup(connection string) (func() error, error) {
	db, err := sql.Open("sqlite3", connection)
	if err != nil {
		return nil, fmt.Errorf("error while connecting to database at %s, %v", connection, err)
	}
	ds.DB = db
	statements := []string{
		"CREATE TABLE IF NOT EXISTS weight (id integer not null primary key, created date not null, weight float(2) not null)",
		"CREATE TABLE IF NOT EXISTS entry (id integer not null primary key, created date not null, entrydate varchar(10) not null, calories integer not null, food varchar(30), bmr float(2) not null, amr float(2) not null)",
		"CREATE TABLE IF NOT EXISTS config (id integer not null primary key, height float(2) not null, activity float(2) not null, birthday date not null, gender varchar(20) not null, unitsystem varchar(20) not null)",
	}
	_, err = db.Exec(strings.Join(statements, ";"))
	if err != nil {
		closeErr := db.Close()
		if closeErr != nil {
			return nil, err
		}
		return nil, err
	}
	return db.Close, nil
}

// SetConfig overrides the current config with the given values
// and deletes the old config and adds a new one
func (ds *SQLiteDataSource) SetConfig(c *model.Config) error {
	_, err := ds.DB.Exec("DELETE from config")
	if err != nil {
		return fmt.Errorf("could not update config: %v", err)
	}
	if c.UnitSystem != "metric" && c.UnitSystem != "imperial" {
		return fmt.Errorf("unit system needs to be either metric or imperial: %s", c.UnitSystem)
	}
	height := c.Height
	if c.UnitSystem == util.Imperial {
		height = util.ToCm(height)
	}
	_, err = ds.DB.Exec("INSERT INTO config (id, height, activity, birthday, gender, unitsystem) VALUES(?, ?, ?, ?, ?, ?)", nil, height, c.Activity, c.Birthday, c.Gender, c.UnitSystem)
	return err
}

// FetchConfig fetches and returns the current config
func (ds *SQLiteDataSource) FetchConfig() (config *model.Config, err error) {
	rows, err := ds.DB.Query("SELECT id, height, activity, birthday, gender, unitsystem FROM config LIMIT 1")
	if err != nil {
		return nil, fmt.Errorf("could not retrieve config: %v", err)
	}
	defer func() {
		closeErr := rows.Close()
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()
	for rows.Next() {
		var id int
		var height float64
		var activity float64
		var birthday time.Time
		var gender string
		var unit string
		err = rows.Scan(&id, &height, &activity, &birthday, &gender, &unit)
		if err != nil {
			return nil, err
		}
		return &model.Config{ID: id, Height: height, Activity: activity, Birthday: birthday, Gender: gender, UnitSystem: unit}, nil
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return nil, errors.New("could not set config")
}

// AddWeight adds the given weight for todays date
func (ds *SQLiteDataSource) AddWeight(weight float64) error {
	config, err := ds.FetchConfig()
	if err != nil {
		return err
	}
	if config.UnitSystem == util.Imperial {
		weight = util.ToKg(weight)
	}
	_, err = ds.DB.Exec("INSERT INTO weight(id, created, weight) VALUES (?, ?, ?)", nil, time.Now(), weight)
	return err
}

// CurrentWeight fetches and returns the current weight, which is the last entry in the table
func (ds *SQLiteDataSource) CurrentWeight() (*model.Weight, error) {
	rows, err := ds.DB.Query("SELECT created, weight FROM weight ORDER BY created DESC LIMIT 1")
	if err != nil {
		return nil, fmt.Errorf("could not fetch current weight: %v", err)
	}
	defer func() {
		closeErr := rows.Close()
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()
	for rows.Next() {
		var created time.Time
		var weight float64
		err = rows.Scan(&created, &weight)
		if err != nil {
			return nil, err
		}
		return &model.Weight{Created: created, Weight: weight}, nil
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return nil, errors.New("weight was not set")
}

// FetchWeights fetches all weight entries
func (ds *SQLiteDataSource) FetchWeights() ([]*model.Weight, error) {
	var result []*model.Weight
	rows, err := ds.DB.Query("SELECT id, created, weight FROM weight ORDER BY created ASC")
	if err != nil {
		return nil, fmt.Errorf("could not retrieve weight history: %v", err)
	}
	defer func() {
		closeErr := rows.Close()
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()
	for rows.Next() {
		var id int
		var created time.Time
		var weight float64

		err = rows.Scan(&id, &created, &weight)
		if err != nil {
			return nil, err
		}

		result = append(result, &model.Weight{ID: id, Created: created, Weight: weight})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, err
}

// AddEntry fetches the current config and weight to calculate the metabolic rate and adds the data
// into the entry table
func (ds *SQLiteDataSource) AddEntry(entryDate string, calories int, food string) error {
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
	_, err = ds.DB.Exec("INSERT INTO entry (id, created, entrydate, calories, food, amr, bmr) VALUES (?, ?, ?, ?, ?, ?, ?)", nil, time.Now(), entryDate, calories, food, amr, bmr)
	if err != nil {
		return fmt.Errorf("could not add entry: %v", err)
	}
	return nil
}

// FetchEntries fetches and returns all entries for a given date
func (ds *SQLiteDataSource) FetchEntries(entryDate string) (model.Entries, error) {
	var res model.Entries
	rows, err := ds.DB.Query("SELECT id, created, entrydate, calories, food, amr, bmr FROM entry WHERE entrydate = ?", entryDate)
	if err != nil {
		return nil, fmt.Errorf("could not fetch entries for the given date: %s, %v", entryDate, err)
	}
	defer func() {
		closeErr := rows.Close()
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()
	res, err = entriesFromRows(rows)
	if err != nil {
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// FetchAllEntries fetches and returns all entries
func (ds *SQLiteDataSource) FetchAllEntries() (model.Entries, error) {
	var res model.Entries
	rows, err := ds.DB.Query("SELECT id, created, entrydate, calories, food, amr, bmr FROM entry")
	if err != nil {
		return nil, fmt.Errorf("could not fetch all entriess, %v", err)
	}
	defer func() {
		closeErr := rows.Close()
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()
	res, err = entriesFromRows(rows)
	if err != nil {
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// entriesFromRows scans the database values of an entry into a result array of model.Entry
func entriesFromRows(rows *sql.Rows) (model.Entries, error) {
	var res model.Entries
	for rows.Next() {
		var id int
		var created time.Time
		var date string
		var calories int
		var food string
		var amr float64
		var bmr float64
		err := rows.Scan(&id, &created, &date, &calories, &food, &amr, &bmr)
		if err != nil {
			return nil, err
		}
		res = append(res, &model.Entry{
			ID:        id,
			Created:   created,
			EntryDate: date,
			Calories:  calories,
			Food:      food,
			AMR:       amr,
			BMR:       bmr,
		})
	}
	return res, nil
}

// RemoveEntries removes all entries for a given day from the database
func (ds *SQLiteDataSource) RemoveEntries(entryDate string) error {
	_, err := ds.DB.Exec("DELETE FROM entry WHERE entrydate = ?", entryDate)
	if err != nil {
		return fmt.Errorf("could not delete entries for %s", entryDate)
	}
	return nil
}

// RemoveEntry removes the entry with the given id for a given day from the database
func (ds *SQLiteDataSource) RemoveEntry(entryDate string, id int) error {
	_, err := ds.DB.Exec("DELETE FROM entry WHERE entrydate = ? AND id = ?", entryDate, id)
	if err != nil {
		return fmt.Errorf("could not delete entry with id %d on day %s", id, entryDate)
	}
	return nil
}

// Import imports the given data to the database, overwriting the previous
// data
func (ds *SQLiteDataSource) Import(data *model.ImpEx) error {
	err := ds.SetConfig(data.Config)
	if err != nil {
		return fmt.Errorf("could not replace config, %v", err)
	}
	for _, weight := range data.Weights {
		_, err = ds.DB.Exec("INSERT OR REPLACE INTO weight (id, created, weight) VALUES (?, ?, ?)", weight.ID, weight.Created, weight.Weight)
		if err != nil {
			return fmt.Errorf("could not insert/update weight with id %d", weight.ID)
		}
	}
	for _, entry := range data.Entries {
		_, err = ds.DB.Exec("INSERT OR REPLACE INTO entry (id, created, entrydate, calories, food, amr, bmr) VALUES (?, ?, ?, ?, ?, ?, ?)", entry.ID, entry.Created, entry.EntryDate, entry.Calories, entry.Food, entry.AMR, entry.BMR)
		if err != nil {
			return fmt.Errorf("could not insert/update entry with id %d", entry.ID)
		}
	}
	return nil
}

// Export creates a JSON representation of the database
func (ds *SQLiteDataSource) Export() (*model.ImpEx, error) {
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
