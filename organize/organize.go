package organize

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Episode stores extracted episode filename metadata
type Episode struct {
	Path     string  // The full path to the video file
	Filename string  // The filename
	Series   string  // The name of the show
	Season   string  // The season as a non-zero padded string
}

// MoveEpisode to the destination subfolder, if it exists
func MoveEpisode(episode Episode, dest string) {
	r := strings.NewReplacer("{Series}", episode.Series,
		"{Season}", episode.Season,
		"{Filename}", episode.Filename)
	destPath := r.Replace(dest)
	destDir := filepath.Dir(destPath)

	if _, err := os.Stat(destDir); !os.IsNotExist(err) {
		err := os.Rename(episode.Path, destPath)
		if err != nil {
			log.Printf("Failed to move %s to %s: %v\n", episode.Path, destPath, err)
		} else {
			log.Printf("Moved %s to %s\n", episode.Path, destPath)
		}
	} else {
		log.Printf("No destination found for %s\n", destDir)
	}
}

// ParseEpisode turns file path into an Episode struct, or fails and returns an error
func ParseEpisode(path string, aliases map[string]string) (Episode, error) {
	var e Episode
	var err error
	e.Path = path
	e.Filename = filepath.Base(path)

	// Pull out episode name and season (stripping off leading zeros)
	r, err := regexp.Compile(`(.+)[sS]0*([\d]+)[eE][\d]+.+`)
	if err != nil {
		return e, err
	}
	res := r.FindStringSubmatch(e.Filename)
	if len(res) != 3 {
		return e, errors.New("Failed to find episode name and season")
	}

	// "." is space, and ToTitle makes each word upper-case-first
	e.Series = strings.TrimSpace(strings.Title(strings.ToLower(strings.ReplaceAll(res[1], ".", " "))))
	e.Season = res[2]

	// Handle episode alises
	if val, ok := aliases[e.Series]; ok {
		// log.Printf("Series alias: `%s` -> `%s`\n", e.Series, val)
		e.Series = val
	}

	return e, err
}

// GetVideoFiles finds all video files that meet the min-size and extension criteria
func GetVideoFiles(folder string, minSize int64, extensions []string) []string {
	var videoFiles []string
	err := filepath.Walk(folder,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Size() >= minSize {
				//fmt.Println(path, info.Size())
				for _, extension := range extensions {
					if strings.HasSuffix(path, extension) {
						// Record file path for later return
						videoFiles = append(videoFiles, path)
						return nil // Done processing this file
					}
				}
				log.Println("Unknown possible video format: %s", path)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return videoFiles
}
