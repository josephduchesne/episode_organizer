package config_test

import (
	"fmt"
	"github.com/josephduchesne/episode_organizer/config"
	"reflect"
	"testing"
)

// TestConfig is the GetConfig Config loading unit test
func TestGetConfig(t *testing.T) {
	var input config.Config
	input.GetConfig("testdata/config.yaml")

	expected := config.Config{
		Source:        "/some/source",
		Extensions:    []string{"A", "B"},
		Aliases:       map[string]string{"Foo": "Bar", "Baz": "Bash"},
		Dest:          "/some/dest",
		MinSize:       123,
		CreateSeasons: true,
	}
	if reflect.DeepEqual(input, expected) {
		fmt.Println("Configuration loaded as expected")
	} else {
		fmt.Printf("Config mismatch:\nExpected: %v\nRead: %v", expected, input)
		t.Fail()
	}
}

// TestMissingConfig tests that an error is thrown when a missing
// config file is loaded
func TestMissingConfig(t *testing.T) {
	var input config.Config
	err := input.GetConfig("testdata/missing_config.yaml")
	if err == nil {
		t.Fail()
	}
}

// TestBadConfig tests that if a config is invalid, it errors
func TestBadConfig(t *testing.T) {
	var input config.Config
	err := input.GetConfig("testdata/invalid_config.yaml")
	if err == nil {
		t.Fail()
	}
}
