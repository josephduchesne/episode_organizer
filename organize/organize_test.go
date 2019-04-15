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

// runEndToEnd is a test runner that validates general functionality of episode_organizer
// It copies intput/output directories in testdata to tmp_input/tmp_output, runs
// and then compares the result

func runEndToEnd(t *testing.T, configFile string) {
	// End-to-end functional test of the main program

	// reset test and copy in new test data
	exec.Command("rm", "-rf", "testdata/tmp_input", "testdata/tmp_output").Run() // Clear old data
	exec.Command("cp", "-a", "testdata/input/", "testdata/tmp_input").Run()
	exec.Command("cp", "-a", "testdata/output/", "testdata/tmp_output").Run()

	// Run the main program
	organize.Episodes(configFile)

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

// End-to-end test without folder creation
func TestMain(t *testing.T) {
	runEndToEnd(t, "testdata/config.yaml")
}

// End-to-end test without folder creation
func TestMainCreate(t *testing.T) {
	runEndToEnd(t, "testdata/config.create.yaml")
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

var episodeTests = []struct {
	in  string
	alias map[string]string
	err bool
	out organize.Episode
}{
	{
		"/foo/bar/show.S01E02.mkv",
		map[string]string{},
		false,
		organize.Episode{"/foo/bar/show.S01E02.mkv", "show.S01E02.mkv", "Show", "1", "2"},
	},
	{
		"/foo/bar/SHOW.s99E102.mkv",
		map[string]string{},
		false,
		organize.Episode{"/foo/bar/SHOW.s99E102.mkv", "SHOW.s99E102.mkv", "Show", "99", "102"},
	},
	{
		"/foo/bar/SHOWwithinvalidname.mkv",
		map[string]string{},
		true,
		organize.Episode{},
	},
}

// TestParseEpisode walks through episodeTests, validing each parse
func TestParseEpisode(t *testing.T) {
	for _, test := range episodeTests {
		episode, err := organize.ParseEpisode(test.in, test.alias)
		if test.err {  // test should fail
			if err == nil {  // we should have errored but didn't
				log.Printf("ParseEpisode Should have errored: %v but got %v,%v\n", test, episode, err)
				t.Fail()
			}
			// Implicitly pass, we failed as expected
			log.Printf("ParseEpisode errored as anticipated for %v: %v\n", test, err)
		} else {  // Regular test, should succeed
			if err != nil { //error
				log.Printf("ParseEpisode Should not have errored: %v but got %v,%v\n", test, episode, err)
				t.Fail()
			}
			if episode != episode {  // mismatch
				log.Printf("ParseEpisode Should match: %v but got %v,%v\n", test, episode, err)
				t.Fail()
			}
			log.Printf("ParseEpisode Matched Expectations: %v and got %v,%v\n", test, episode, err)
		}
	}
}
