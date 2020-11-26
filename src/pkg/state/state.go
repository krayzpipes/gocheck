// Package state provides a way for scheduled jobs to
// persist state to
package state

import (
	"io/ioutil"
	"os"
	"strings"
)

// Variadic, write to file on disk.
	// Newline delimited
// Read from disk

// StateField holds the key/value pair that will
// be persisted to file.
type StateField struct {
	Name string
	Value string
}

func (s *StateField) FileFmtString() string {
	return s.Name + "=" + s.Value
}

// State will contain the location the state
// will be stored as well as the fields that will
// be stored in the state file
type State struct {
	Location string
	Fields map[string]*StateField
}

func (s *State) FieldsToString() string {
	var fields []string
	for field := range s.Fields {
		fieldString := s.Fields[field].FileFmtString()
		fields = append(fields, fieldString)
	}
	return strings.Join(fields, "\n")
}

func (s *State) Save() error {
	f, err := os.Create(s.Location)  // Overwrite everytime we save
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(s.FieldsToString())
	if err != nil {
		return err
	}

	f.Sync()
	return nil
}

func (s *State) Load() error {
	data, err := ioutil.ReadFile(s.Location)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")

	// TODO - What if keys already exist?... do we need to
	// check for that here?
	for _, line := range lines {
		fieldPair := strings.Split(line, "=")
		s.Fields[fieldPair[0]] = &StateField{Name: fieldPair[0], Value: fieldPair[1]}
	}
	return nil
}