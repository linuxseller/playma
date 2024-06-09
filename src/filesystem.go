package main
import (
    "os"
    "strings"
    "log"
    "os/exec"
    "bytes"
    "strconv"
    "encoding/csv"
    "errors"
    "path/filepath"
    "slices"
)
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

const music_data_file = "playma_data/songs_db.csv"

func ParseSongDataFromFilename(name string) (title, artist string) {
    splitted := strings.Split(TrimFileExtention(name), " - ")
    if len(splitted)>1{
        title = strings.Trim(splitted[1], " ")
        artist = strings.Trim(splitted[0], " ")
    } else {
        title = TrimFileExtention(name)
        artist = "unknown"
    }
    log.Print(title, artist)
    return title, artist
}

func TraverseMusicDir(dir, path string) ([]Song, error) {
    var res []Song
    dir_content, err := os.ReadDir(dir)
    if err != nil {
        return nil, err
    }
    for _, val := range dir_content {
        path_to_file := dir + "/" + val.Name()
        if val.IsDir() {
            deeper, err := TraverseMusicDir(path_to_file, dir)
            if err != nil {
                return nil, err
            }
            res = append(res, deeper...)
        } else {
            if !slices.Contains(supported_filetypes, filepath.Ext(val.Name())){
                continue
            }
            cmd := exec.Command("ffmpeg", "-i", path_to_file, "-f", "ffmetadata", "pipe:1")
            var out bytes.Buffer
            cmd.Stdout = &out
            err := cmd.Run()
            keyvals, err := ParseKeyVal(out.String())
            if err != nil {
                log.Fatal(err)
            }
            var song Song
            title, artist := ParseSongDataFromFilename(val.Name())
            if _, ok := keyvals["title"]; ok {
                song.Title = keyvals["title"]
            } else {song.Title = title}

            if _, ok := keyvals["artist"]; ok {
                song.Artist = keyvals["artist"]
            } else {song.Artist = artist}

            if _, ok := keyvals["album"]; ok {
                song.Album = keyvals["album"]
            } else {song.Album = "undefined"}

            if _, ok := keyvals["track"]; ok {
                song.Track_num,err = strconv.Atoi(keyvals["track"])
                if err != nil {
                    song.Track_num = 0;
                }
            } else {song.Track_num = 0}

            if _, ok := keyvals["genre"]; ok {
                song.Genre = keyvals["genre"]
            } else {song.Genre = "undefined"}

            if _, ok := keyvals["album_artist"]; ok {
                song.Album_artist = keyvals["album_artist"]
            } else {song.Album_artist = "undefined"}

            song.Path = path_to_file
            res = append(res, song)
        }
    }
    return res, nil
}
func SaveSongs(songs []Song) error {
    file, err := os.Create(music_data_file)
    if err != nil {
        return err
    }
    defer file.Close()
    writer := csv.NewWriter(file)
    defer writer.Flush()
    for _, i := range songs {
        writer.Write(i.ToSlice())
    }
    return nil
}
func LoadSongs() ([]Song, error) {
    var songs []Song
    if !fileExists(music_data_file) {
        return nil, errors.New("could not open " + music_data_file)
    }
    file, err := os.Open(music_data_file)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    reader := csv.NewReader(file)
    loaded, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }
    for _, i := range loaded {
        songs = append(songs, SliceToSong(i))
    }
    return songs, nil
}

func UpdateSongs() ([]Song, error) {
    log.Println("started tree traverse")
    var new_music_files []Song
    for _, i := range config.Music_dirs {
        i_music_files, err := TraverseMusicDir(i, i)
        if err!=nil {
            log.Println("Error while traversing music dir", err)
        }
        new_music_files = append(new_music_files, i_music_files...)
    }
    log.Printf("finished tree traverse of %d files\n", len(new_music_files))
    err := SaveSongs(new_music_files)
    if err != nil {
            log.Fatal("Error: could not save songs to db", err)
    }
    return new_music_files, nil
}

