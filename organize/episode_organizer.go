package organize

import (
	"github.com/josephduchesne/episode_organizer/config"
	"log"
)

// Episodes is the main program that loads config.yaml,
// finds video files and moves them into target folders
func Episodes(configFile string) {
	var c config.Config
	c.GetConfig(configFile)

	log.Printf("Config: %+v\n\n", c)

	videoFiles := GetVideoFiles(c.Source, c.MinSize, c.Extensions)
	for _, file := range videoFiles {
		episode, err := ParseEpisode(file, c.Aliases)
		if err != nil {
			log.Printf("Error parsing episode %s: %v", file, err)
		} else {
			MoveEpisode(episode, c.Dest)
		}
	}
}
