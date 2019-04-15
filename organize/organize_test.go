package organize_test

import (
	"bytes"
	"flag"
	"github.com/josephduchesne/episode_organizer/organize"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"testing"
)

// To update the test's expected output, run `go test -update`
var update = flag.Bool("update", false, "update .golden files")

// TestMain is an end-to-end test that validates general functionality of episode_organizer
// It copies intput/output directories in testdata to tmp_input/tmp_output, runs
// and then compares the result
func TestMain(t *testing.T) {
	// End-to-end functional test of the main program

	// reset test and copy in new test data
	exec.Command("rm", "-rf", "testdata/tmp_input", "testdata/tmp_output").Run() // Clear old data
	exec.Command("cp", "-a", "testdata/input/", "testdata/tmp_input").Run()
	exec.Command("cp", "-a", "testdata/output/", "testdata/tmp_output").Run()

	// Run the main program
	organize.Episodes("testdata/config.yaml")

	// Check the results (with optional update flag)
	actual, err := exec.Command("sh", "-c", "find testdata/tmp_input testdata/tmp_output | sort").CombinedOutput()
	if err != nil {
		log.Printf("Error inspecting test output: %v\n\t%s\n", err, actual)
		t.Fail()
	}
	golden := filepath.Join("testdata", t.Name()+".golden")
	if *update {
		ioutil.WriteFile(golden, actual, 0644)
	}
	expected, _ := ioutil.ReadFile(golden)

	if !bytes.Equal(actual, expected) {
		log.Println("Output was not what we expected! Please see the resulting directory structure.")
		log.Printf("Expected:\n%s\n\n", expected)
		log.Printf("Actual:\n%s\n\n", actual)
		t.Fail()
	}
}

// TestMainFailConfig checks the case where config fails to load
func TestMainFailConfig(t *testing.T) {
	err := organize.Episodes("testdata/invalid_config.yaml")
	if err == nil {
		log.Println("Config load should have errored")
		t.Fail()
	}
}

// TestFailedRename checks the case where an episode can't be moved into place
func TestFailedRename(t *testing.T) {
	organize.Episodes("testdata/config.yaml")
	// Then try to move one again
	episode := organize.Episode{Path: "/tmp/a/very/fake/path", Filename: "Foo", Series: "Real", Season: "1"}
	err := organize.MoveEpisode(episode, "testdata/tmp_output/")
	if err == nil {
		log.Println("TestFailedRename should have errored")
		t.Fail()
	}
}

// TestNoDest checks the case where the target directory doesn't exist
func TestNoDest(t *testing.T) {
	organize.Episodes("testdata/config.yaml")
	// Then try to move one again
	episode := organize.Episode{Path: "testdata/tmp_input/No_Dest.S07E02.mkv", Filename: "Foo", Series: "Real", Season: "1"}
	err := organize.MoveEpisode(episode, "/dev/null/fake")
	if err == nil {
		log.Println("MoveEpisode should have errored")
		t.Fail()
	}
}

// TestBadSourceFolder tests that errors are handled properly when the source folder is invalid
func TestBadSourceFolder(t *testing.T) {
	_, err := organize.GetVideoFiles("/dev/null/fake", 0, []string{"foo"})
	if err == nil {
		log.Println("GetVideoFiles should have an error")
		t.Fail()
	}
}

// TestBadSourceConfig tests error handling for a bad source file in a configuration
func TestBadSourceConfig(t *testing.T) {
	err := organize.Episodes("testdata/config_bad_source.yaml")
	if err == nil {
		log.Println("Episode organization should have failed with a bad source dir")
		t.Fail()
	}
}
