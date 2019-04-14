package organize

import (
    "os"
    "log"
    "strings"
    "path/filepath"
    "regexp"
    "errors"
)

type Episode struct {
    Path string
    Filename string
    Series string
    Season string
}

// Move episode to the destination subfolder, if it exists
func MoveEpisode(episode Episode, dest string) {
    r := strings.NewReplacer("{Series}", episode.Series,
                             "{Season}", episode.Season,
                             "{Filename}", episode.Filename)
    dest_path := r.Replace(dest)
    dest_dir  := filepath.Dir(dest_path)

    if _, err := os.Stat(dest_dir); !os.IsNotExist(err) {
        err := os.Rename(episode.Path, dest_path)
        if err != nil {
            log.Printf("Failed to move %s to %s: %v\n", episode.Path, dest_path, err)
        } else {
            log.Printf("Moved %s to %s\n", episode.Path, dest_path)
        }
    } else {
        log.Printf("No destination found for %s\n", dest_dir)
    }
}

func ParseEpisode(path string, aliases map[string]string) (Episode, error) {
    var e Episode
    var err error
    e.Path = path
    e.Filename = filepath.Base(path)

    // Pull out episode name and season (stripping off leading zeros)
    r, err := regexp.Compile(`(.+)[sS]0*([\d]+)[eE][\d]+.+`)
    if err != nil {
        return e,err
    }
    res := r.FindStringSubmatch(e.Filename)
    if len(res) != 3 {
        return e,errors.New("Failed to find episode name and season")
    }

    // "." is space, and ToTitle makes each word upper-case-first
    e.Series = strings.Title(strings.ToLower(strings.ReplaceAll(res[1], ".", " ")))
    e.Season = res[2]

    // Handle episode alises
    if val, ok := aliases[e.Series]; ok {
        // log.Printf("Series alias: `%s` -> `%s`\n", e.Series, val)
        e.Series = val
    }

    return e, err
}

// Find all video files that meet the min-size and extension criteria
func GetVideoFiles(folder string, min_size int64, extensions []string) []string {
    var video_files []string;
    err := filepath.Walk(folder,
        func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if (info.Size()>=min_size){
            //fmt.Println(path, info.Size())
            for _, extension := range extensions {
                if strings.HasSuffix(path, extension) {
                    // Record file path for later return
                    video_files = append(video_files, path)
                    return nil  // Done processing this file
                }
            }
            log.Println("Unknown possible video format: %s", path)
        }
        return nil
    })
    if err != nil {
        log.Println(err)
    }
    return video_files
}
