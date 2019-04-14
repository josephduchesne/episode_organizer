package main

import (
	"github.com/josephduchesne/episode_organizer/config"
	"github.com/josephduchesne/episode_organizer/organize"
	"log"
)

func main() {
	var c config.Config
	c.GetConfig()

	log.Printf("Config: %+v\n\n", c)

	video_files := organize.GetVideoFiles(c.Source, c.MinSize, c.Extensions)
	for _, file := range video_files {
		episode, err := organize.ParseEpisode(file, c.Aliases)
		if err != nil {
			log.Printf("Error parsing episode %s: %v", file, err)
		} else {
			organize.MoveEpisode(episode, c.Dest)
		}
	}
}
