package main

import (
	"github.com/josephduchesne/episode_organizer/config"
	"github.com/josephduchesne/episode_organizer/organize"
	"log"
)

func main() {
	var c config.Config
	c.GetConfig("config.yaml")

	log.Printf("Config: %+v\n\n", c)

	videoFiles := organize.GetVideoFiles(c.Source, c.MinSize, c.Extensions)
	for _, file := range videoFiles {
		episode, err := organize.ParseEpisode(file, c.Aliases)
		if err != nil {
			log.Printf("Error parsing episode %s: %v", file, err)
		} else {
			organize.MoveEpisode(episode, c.Dest)
		}
	}
}
