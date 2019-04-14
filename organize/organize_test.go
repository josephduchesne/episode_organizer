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
