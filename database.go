package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Database struct {
	databasePath string
	CurrentId    uint     `json:"currentId"`
	Entires      []*Entry `json:"entries"`
}

func CreateDatabase(path string, data []byte) (*Database, error) {
	var db Database
	db.databasePath = path
	if len(data) == 0 {
		db.CurrentId = 1
		db.Entires = []*Entry{}
		return &db, nil
	}

	err := json.Unmarshal([]byte(data), &db)
	if err != nil {
		return nil, fmt.Errorf("error parsing database %w", err)
	}

	fmt.Println("db path", db.databasePath)

	return &db, nil
}

func (d *Database) Add(name string, start, end time.Time) (*Entry, error) {
	entry := &Entry{
		Id:    d.CurrentId,
		Name:  name,
		Start: start,
		End:   end,
	}
	err := validateEntry(*entry)
	if err != nil {
		return nil, err
	}
	d.CurrentId += 1

	d.Entires = append(d.Entires, entry)

	err = d.save()
	if err != nil {
		return nil, fmt.Errorf("error adding: %w", err)
	}

	return entry, nil
}

func (d *Database) GetAll() []*Entry {
	return d.Entires
}

func (d *Database) GetById(id uint) *Entry {
	entry, _, _ := d.getEntryAndIndexById(id)
	return entry
}

func (d *Database) getEntryAndIndexById(id uint) (*Entry, uint, error) {
	var idx uint
	var entry *Entry
	for i, e := range d.Entires {
		if e.Id == id {
			idx = uint(i)
			entry = e
			break
		}
	}
	if entry == nil {
		return nil, 0, errors.New("entry not found")
	}
	return entry, idx, nil
}

func (d *Database) DeleteById(id uint) (*Entry, error) {
	entry, idx, err := d.getEntryAndIndexById(id)
	if err != nil {
		return nil, err
	}
	d.Entires = append(d.Entires[:idx], d.Entires[idx+1:]...)

	err = d.save()
	if err != nil {
		return nil, fmt.Errorf("error deleting: %w", err)
	}

	return entry, nil
}

func (d *Database) UpdateById(id uint, toUpdate Entry) (*Entry, error) {
	err := validateEntry(toUpdate)
	if err != nil {
		return nil, err
	}
	entry, _, err := d.getEntryAndIndexById(id)
	if err != nil {
		return nil, err
	}
	entry.Start = toUpdate.Start
	entry.End = toUpdate.End
	entry.Name = toUpdate.Name

	err = d.save()
	if err != nil {
		return nil, fmt.Errorf("error updating: %w", err)
	}

	return entry, nil
}

func (d *Database) save() error {
	data, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("error marshalling database %w", err)
	}

	err = writeFileContents(d.databasePath, data)
	if err != nil {
		return fmt.Errorf("error saving database %w", err)
	}

	return nil
}
