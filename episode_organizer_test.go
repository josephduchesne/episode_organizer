package main

import (
	"bytes"
	"flag"
	"github.com/josephduchesne/episode_organizer/config"
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
// and then comapres the result
func TestMain(t *testing.T) {
	// End-to-end functional test of the main program
	var c config.Config
	c.GetConfig("testdata/config.yaml")

	log.Printf("Config: %+v\n\n", c)

	// reset test
	exec.Command("rm", "-rf", "testdata/tmp_input", "testdata/tmp_output").Run() // Clear old data
	// Then copy in the new test data
	exec.Command("cp", "-a", "testdata/input/", "testdata/tmp_input").Run()
	exec.Command("cp", "-a", "testdata/output/", "testdata/tmp_output").Run()

	videoFiles := organize.GetVideoFiles(c.Source, c.MinSize, c.Extensions)
	for _, file := range videoFiles {
		episode, err := organize.ParseEpisode(file, c.Aliases)
		if err != nil {
			log.Printf("Error parsing episode %s: %v", file, err)
		} else {
			organize.MoveEpisode(episode, c.Dest)
		}
	}

	// Check the results (with optional update flag)
	actual, err := exec.Command("find", "testdata/tmp_input", "testdata/tmp_output").CombinedOutput()
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
		t.Fail()
	}
}
