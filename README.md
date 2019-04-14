# What is this?

This is a small tool to make TV show episode organization easier without risking really messing up your directory structure, file names, or other things.

It will find all files above a specified size that has specified extensions and move them to a specified destination folder, organized by Series name and Episode.

The main design principle is to do nothing if anything not exactly as expected. 

The means that if something isn't handled properly, you can ammend the program or configuration to resolve the issue and rerun the program.

This program expects that the series and season folders will exist already. It would be pretty easy to add optional Series or Season folder creation.

# How to install

- [Install Go](https://golang.org/doc/install)
- Download source with `go get github.com/josephduchesne/episode_organizer`
- Build with `go build github.com/josephduchesne/episode_organizer`
- Create a config file `wget https://raw.githubusercontent.com/josephduchesne/episode_organizer/master/config.yaml.dist -O config.yaml`
- Edit config.yaml
- Run with `./episode_organizer`

## Spiffy Badges
[![Go Report Card](https://goreportcard.com/badge/github.com/josephduchesne/episode_organizer)](https://goreportcard.com/report/github.com/josephduchesne/episode_organizer) 
