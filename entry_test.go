package main

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateEntry(t *testing.T) {
	testStart := time.Now().Add(-time.Second)
	testEnd := time.Now()
	tests := []struct {
		name          string
		inputEntry    Entry
		expectedError error
	}{
		{
			name: "happy path - valid entry",
			inputEntry: Entry{
				Id:    1,
				Name:  "streaming",
				Start: testStart,
				End:   testEnd,
			},
		},
		{
			name: "unhappy path - start is after end",
			inputEntry: Entry{
				Id:    1,
				Name:  "streaming",
				Start: testEnd,
				End:   testStart,
			},
			expectedError: errors.New("start must be before end"),
		},
		{
			name: "unhappy path - empty name",
			inputEntry: Entry{
				Id:    1,
				Name:  "",
				Start: testStart,
				End:   testEnd,
			},
			expectedError: errors.New("name must not be empty"),
		},
	}
	for _, ttest := range tests {
		t.Run(ttest.name, func(tt *testing.T) {
			err := validateEntry(ttest.inputEntry)
			assert.Equal(tt, ttest.expectedError, err)
		})
	}
}
