package main

import (
    "fmt"
    "strings"
    "log"
    "net/http"
    "io/ioutil"
    // "path/filepath"
)

const config_file = "playma_data/config.txt"

var supported_filetypes []string = []string{".mp3", ".wav", ".opus"}

type Config struct {
    Music_dirs []string
    Playlist_dir string
}

type Song struct {
    Title        string
    Artist       string
    Album        string
    Track_num    int
    Genre        string
    Album_artist string
    Path         string
}
func (song Song) GetHeaders() []string {
    return []string{"Title","Artist","Album","Track_num","Genre","lbum_artist","Path"}
}
func (song Song) ToSlice() []string {
    return []string{song.Artist, song.Title, song.Album, /*song.Track_num,*/ song.Genre, song.Album_artist, song.Path}
}

func SliceToSong(slice []string) Song {
    var song Song
    song.Artist       = slice[0]
    song.Title        = slice[1]
    song.Album        = slice[2]
    song.Genre        = slice[3]
    song.Album_artist = slice[4]
    song.Path         = slice[5]
    return song
}

func (a Song)Equal(b Song) bool {
    return a.Title == b.Title && a.Artist == b.Title
}

type HomeData struct {
    Songs []Song
}

func TrimFileExtention(name string) (result string) {
    for _, i := range supported_filetypes {
        result = strings.TrimSuffix(name, i)
        if result != name {
            return result
        }
    }
    return name
}

var config Config
var music_files []Song
func main() {
    file, err := ioutil.ReadFile(config_file)
    if err != nil {
        log.Println("Error reading config file:", err)
        return
    }
    read_config, err := ParseKeyVal(string(file))
    if err != nil {
        log.Fatal(err)
        return
    }
    config.Music_dirs = strings.Split(strings.Trim(read_config["music_dir"], "\""), ":")
    config.Playlist_dir = read_config["playlist_dir"]
    fmt.Printf("Config: %+v\n", config)
    music_files, err = LoadSongs()
    if err != nil {
        log.Printf("Error loading songs db, will only traverse tree")
    }
    go func() {
        fs := http.FileServer(http.Dir("./"))
        http.Handle("/favicon.ico", fs)
        http.Handle("/images/", fs)
        http.HandleFunc("/music/", song)
        http.HandleFunc("/home", home)
        log.Println("Running server on :8080")
        http.ListenAndServe(":8080", nil)
    }()
    music_files, err = UpdateSongs()
    select {}
}
