package organize

import (
	"github.com/josephduchesne/episode_organizer/config"
	"log"
)

// Episodes is the main program that loads config.yaml,
// finds video files and moves them into target folders
func Episodes(configFile string) error {
	var c config.Config
	configError := c.GetConfig(configFile)
	if configError != nil {
		return configError
	}

	log.Printf("Config: %+v\n\n", c)

	videoFiles, videoFileError := GetVideoFiles(c.Source, c.MinSize, c.Extensions)

	if videoFileError != nil {
		return videoFileError
	}

	for _, file := range videoFiles {
		episode, err := ParseEpisode(file, c.Aliases)
		if err != nil {
			log.Printf("Error parsing episode %s: %v", file, err)
		} else {
			MoveEpisode(episode, c.Dest, c.CreateSeasons)
		}
	}
	return nil
}
