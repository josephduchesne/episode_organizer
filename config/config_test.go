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
		Source:     "/some/source",
		Extensions: []string{"A", "B"},
		Aliases:    map[string]string{"Foo": "Bar", "Baz": "Bash"},
		Dest:       "/some/dest",
		MinSize:    123,
	}
	if reflect.DeepEqual(input, expected) {
		fmt.Println("Configuration loaded as expected")
	} else {
		fmt.Printf("Config mismatch:\nExpected: %v\nRead: %v", expected, input)
		t.Fail()
	}
}
