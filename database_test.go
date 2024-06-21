package main

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createGoldenEntry(id uint, name string) *Entry {
	testStart := time.Now().Add(-time.Second)
	testEnd := time.Now()
	return &Entry{
		Id:    id,
		Name:  name,
		Start: testStart,
		End:   testEnd,
	}
}

func Test_CreateDatabase(t *testing.T) {
	tests := []struct {
		name           string
		data           string
		expectedResult *Database
		expectedError  error
	}{
		{
			name: "happy path: can load an empty database",
			data: "{\"currentId\": 1, \"entries\": []}",
			expectedResult: &Database{
				CurrentId: 1,
				Entires:   []*Entry{},
			},
		},
		{
			name:          "unhappy path: empty string",
			data:          "",
			expectedError: errors.New("error parsing database unexpected end of JSON input"),
		},
		{
			name:          "unhappy path: junk data",
			data:          "adgflkdsfhgsdfg",
			expectedError: errors.New("error parsing database invalid character 'a' looking for beginning of value"),
		},
	}

	for _, ttest := range tests {
		t.Run(ttest.name, func(tt *testing.T) {
			result, err := CreateDatabase("./test.json", []byte(ttest.data))
			if ttest.expectedError != nil {
				assert.EqualError(tt, ttest.expectedError, err.Error())
			} else {
				assert.Nil(tt, err)
			}
			assert.Equal(tt, ttest.expectedResult, result)
		})
	}
}

func TestDatabase_Add(t *testing.T) {
	testStart := time.Now().Add(-time.Second)
	testEnd := time.Now()
	tests := []struct {
		name           string
		inputName      string
		inputStart     time.Time
		inputEnd       time.Time
		expectedResult *Entry
		expectedLength int
		expectedError  error
	}{
		{
			name:           "happy path",
			inputName:      "Streaming",
			inputStart:     testStart,
			inputEnd:       testEnd,
			expectedLength: 1,
			expectedResult: &Entry{
				Id:    1,
				Name:  "Streaming",
				Start: testStart,
				End:   testEnd,
			},
		},
		{
			name:           "unhappy: name cannot be empty",
			inputName:      "",
			inputStart:     testStart,
			inputEnd:       testEnd,
			expectedLength: 0,
			expectedResult: nil,
			expectedError:  errors.New("name must not be empty"),
		},
		{
			name:           "unhappy: start after end",
			inputName:      "streaming",
			inputStart:     testEnd,
			inputEnd:       testStart,
			expectedLength: 0,
			expectedResult: nil,
			expectedError:  errors.New("start must be before end"),
		},
	}
	for _, ttest := range tests {
		t.Run(ttest.name, func(tt *testing.T) {
			db := &Database{
				CurrentId: 1,
				Entires:   []*Entry{},
			}

			result, err := db.Add(ttest.inputName, ttest.inputStart, ttest.inputEnd)
			assert.Equal(tt, ttest.expectedError, err)
			assert.Equal(tt, ttest.expectedResult, result)
			assert.Equal(tt, ttest.expectedLength, len(db.Entires))
		})
	}
}

func TestDatabase_DeleteById(t *testing.T) {
	goldenEntry := createGoldenEntry(1, "streaming")
	goldenEntryTwo := createGoldenEntry(2, "not streaming")
	goldenEntryThree := createGoldenEntry(3, "maybe streaming")

	tests := []struct {
		name           string
		inputId        uint
		database       *Database
		expectedResult *Entry
		expectedLength int
		expectedError  error
	}{
		{
			name:    "happy path: can delete an entry that exists",
			inputId: 2,
			database: &Database{
				Entires: []*Entry{
					goldenEntry,
					goldenEntryTwo,
					goldenEntryThree,
				},
			},
			expectedResult: goldenEntryTwo,
			expectedLength: 2,
		},
		{
			name:    "unhappy path: cannot delete an entry that does not exist",
			inputId: 2,
			database: &Database{
				Entires: []*Entry{},
			},
			expectedLength: 0,
			expectedError:  errors.New("entry not found"),
		},
	}

	for _, ttest := range tests {
		t.Run(ttest.name, func(tt *testing.T) {
			result, err := ttest.database.DeleteById(ttest.inputId)
			assert.Equal(tt, ttest.expectedError, err)
			assert.Equal(tt, ttest.expectedResult, result)
			assert.Equal(tt, ttest.expectedLength, len(ttest.database.Entires))
		})
	}
}

func TestDatabase_GetById(t *testing.T) {
	goldenEntry := createGoldenEntry(1, "streaming")
	tests := []struct {
		name           string
		database       *Database
		idToSearch     uint
		expectedResult *Entry
	}{
		{
			name:       "happy path - can get entry that exists by id",
			idToSearch: 1,
			database: &Database{
				Entires: []*Entry{
					goldenEntry,
				},
			},
			expectedResult: goldenEntry,
		},
		{
			name:       "unhappy path - get nil when entry doesn't exist",
			idToSearch: 1,
			database: &Database{
				Entires: []*Entry{},
			},
		},
	}

	for _, ttest := range tests {
		t.Run(ttest.name, func(tt *testing.T) {
			result := ttest.database.GetById(ttest.idToSearch)
			assert.Equal(tt, ttest.expectedResult, result)
		})
	}
}

func TestDatabase_UpdateById(t *testing.T) {
	toUpdateStart := time.Now().Add(-time.Second * 10)
	toUpdateEnd := time.Now().Add(time.Second * 5)
	goldenEntry := createGoldenEntry(1, "streaming")
	goldenEntryTwo := createGoldenEntry(2, "not streaming")
	goldenEntryThree := createGoldenEntry(3, "maybe streaming")

	expectedUpdatedEntry := Entry{
		Id:    1,
		Name:  "can update entries",
		Start: toUpdateStart,
		End:   toUpdateEnd,
	}

	tests := []struct {
		name           string
		inputId        uint
		inputToUpdate  Entry
		database       *Database
		expectedResult *Entry
		expectedError  error
	}{
		{
			name:          "happy path - can update an existing entry",
			inputId:       1,
			inputToUpdate: expectedUpdatedEntry,
			database: &Database{
				Entires: []*Entry{
					goldenEntry,
					goldenEntryTwo,
					goldenEntryThree,
				},
			},
			expectedResult: &expectedUpdatedEntry,
		},
		{
			name:          "unhappy path - cannot delete entry that doesn't exist",
			inputId:       1,
			inputToUpdate: expectedUpdatedEntry,
			database: &Database{
				Entires: []*Entry{},
			},
			expectedError: errors.New("entry not found"),
		},
		{
			name:    "unhappy path - toUpdate name is empty",
			inputId: 1,
			inputToUpdate: Entry{
				Id:    1,
				Name:  "",
				Start: toUpdateStart,
				End:   toUpdateEnd,
			},
			database: &Database{
				Entires: []*Entry{
					goldenEntry,
					goldenEntryTwo,
					goldenEntryThree,
				},
			},
			expectedError: errors.New("name must not be empty"),
		},
	}
	for _, ttest := range tests {
		t.Run(ttest.name, func(tt *testing.T) {
			result, err := ttest.database.UpdateById(ttest.inputId, ttest.inputToUpdate)
			assert.Equal(tt, ttest.expectedError, err)
			assert.Equal(tt, ttest.expectedResult, result)
		})
	}
}
