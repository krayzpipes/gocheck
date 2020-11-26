package state

import (
	"log"
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
	nameOne := "field1"
	valueOne := "value1"
	nameTwo := "field2"
	valueTwo := "value2"

	expected := nameOne + "=" + valueOne + "\n" + nameTwo + "=" + valueTwo

	stateFields := make(map[string]*StateField)
	stateFields[nameOne] = &StateField{Name: nameOne, Value: valueOne}
	stateFields[nameTwo] = &StateField{Name: nameTwo, Value: valueTwo}

	state := State{File: "/fake/fileName.gocheck.state", Fields: stateFields}

	actual := state.FieldsToString()
	if expected != actual {
		log.Fatalf("\n----Expected----\n%s\n-----Actual-----\n%s\n----------------", expected, actual)
	}
}

func TestState_Save_NewFile(t *testing.T) {

}