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

func (d *Database) GetAll(filter EntryFilter) []*Entry {
	var filtered []*Entry
	// if filter.ProjectName != "" {
	// 	var filtered []*Entry
	// 	for _, entry := range d.Entires {
	// 		if entry.Name == filter.ProjectName {
	// 			filtered = append(filtered, entry)
	// 		}
	// 	}
	// 	return filtered
	// }
	sfProjectName := filter.ProjectName != ""
	sfSpecificDate := filter.SpecificDate != ""
	for _, entry := range d.Entires {
		shouldAdd := true
		if sfProjectName && entry.Name != filter.ProjectName {
			shouldAdd = false
		}
		edate := entry.Start.Format("2006-01-02")
		if sfSpecificDate && edate != filter.SpecificDate {
			shouldAdd = false
		}
		if shouldAdd {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func (d *Database) GetById(id uint) *Entry {
	entry, _, _ := d.getEntryAndIndexById(id)
	return entry
}

func (d *Database) GetUniqueProjects() []*ProjectCount {
	projects := make([]*ProjectCount, 0)

	for _, entry := range d.Entires {
		foundProject := false
		for _, recordedProject := range projects {
			if entry.Name == recordedProject.Name {
				recordedProject.Count += 1
				foundProject = true
				break
			}
		}
		if foundProject {
			continue
		}
		projects = append(projects, &ProjectCount{
			Name:  entry.Name,
			Count: 1,
		})
	}

	return projects
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
