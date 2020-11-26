package state

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStateField_FileFmtString_Success(t *testing.T) {
	name := "test_field"
	value := "test_value"
	field := StateField{Name: name, Value: value}

	formatted := field.FileFmtString()

	expected := name + "=" + value
	if expected != formatted {
		log.Fatalf("expected %q, got %q", expected, formatted)
	}
}

func TestState_FieldsToString_Success(t *testing.T) {
	state, expected := generateStateAndFields()

	actual := state.FieldsToString()
	if !compareStateContent(expected, actual){
		log.Fatalf("\n----Expected----\n%s\n-----Actual-----\n%s\n----------------", expected, actual)
	}
}

func TestState_Save_NewFile(t *testing.T) {
	state, expected := generateStateAndFields()

	err := state.Save()
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(state.Location)

	content, err := ioutil.ReadFile(state.Location)
	if err != nil {
		log.Fatalf("unable to read file: %s", err)
	}

	text := string(content)

	if !compareStateContent(expected, text) {
		log.Fatalf("\n----Expected----\n%s\n-----Actual-----\n%s\n----------------", expected, text)
	}
}

func TestState_Save_OverwriteFile(t *testing.T) {
	state, expected := generateStateAndFields()

	_ = ioutil.WriteFile(state.Location, []byte("nope=nope\nnope=nope\nnada=nada"), 0777)

	defer os.Remove(state.Location)

	state.Save()

	content, err := ioutil.ReadFile(state.Location)
	if err != nil {
		log.Fatal(err)
	}

	text := string(content)

	if !compareStateContent(expected, text) {
		log.Fatalf("\n----Expected----\n%s\n-----Actual-----\n%s\n----------------", expected, text)
	}
}

func TestState_Load_Success(t *testing.T) {
	state, expected := generateStateAndFields()

	_ = state.Save()

	defer os.Remove(state.Location)

	var secondState State

	secondState.Location = state.Location

	secondState.Load()
	actual := secondState.FieldsToString()

	if !compareStateContent(expected, actual) {
		log.Fatalf(
			"states do not match after load:\n" +
				"----Expected----" +
				"\n%s\n-----Actual-----" +
				"\n%s\n----------------",
				expected,
				actual,
			)
	}
}

func generateStateAndFields() (State, string) {
	tempDir := os.TempDir()
	stateLocation := filepath.Join(tempDir, "gocheck-state-test.gocheck")

	nameOne := "field1"
	valueOne := "value1"
	nameTwo := "field2"
	valueTwo := "value2"

	expectedContent := nameOne + "=" + valueOne + "\n" + nameTwo + "=" + valueTwo

	stateFields := make(map[string]*StateField)
	stateFields[nameOne] = &StateField{Name: nameOne, Value: valueOne}
	stateFields[nameTwo] = &StateField{Name: nameTwo, Value: valueTwo}

	state := State{Location: stateLocation, Fields: stateFields}
	return state, expectedContent
}

// Order of fields does not matter, we just
// want to make sure the same fields/values exist.
func compareStateContent(a string, b string) bool {
	firstSlice := strings.Split(a, "\n")
	secondSlice := strings.Split(b, "\n")

	if len(firstSlice) != len(secondSlice) {
		return false
	}
	for _, field := range firstSlice {
		if !findElementInSlice(field, secondSlice) {
			return false
		}
	}
	for _, field := range secondSlice {
		if !findElementInSlice(field, firstSlice) {
			return false
		}
	}
	return true
}

func findElementInSlice(target string, slice []string) bool {
	for _, item := range slice {
		if target == item {
			return true
		}
	}
	return false
}