# PLAYMA
aka. Play My Audio is a self-hostable simple music server.

# Depemdecies
- go compiler
- ffmpeg, for metadata extraction

## Configuration

config files are in simple ini-like format, `key=value`. Lines starting from `;` are considered comments. Comments after data are forbidden

```
a=b ;this comment is forbidden
```

```
;playma_data/config.txt
music_dir="/home/user/Music:/media/hdd/Music"
playlist_dir="None"
```

## How does this work

When you run `src/playma` (or `./run.sh`) playma first opens `playma_data/confiog.txt` and traverses music folder for `.mp3` files, also it uses ffmpeg to get metadata (if metadata is not present song is read in format `<artist> - <songname>.mp3` spaces are significant).
After list of songs is created it is saved to `playma_data/songs_db.csv`


## Project file structure of modifyable files

You can easily modify some aspects of how playma looks.

```
.
├── favicon.ico
├── images
│   ├── ball.png
│   ├── muted.png
│   ├── pause.png
│   ├── play.png
│   └── volume.png
├── playma_data
│   └── config.txt
└── templates
    └── home.html
```

* A least for now `images` folder contains only images for audio control.
* `home.html` contains css js and html for displaying everything
* `config.txt` contains directories for where to search data

# Project "roadmap"
 - [ ] Support for `.m3u` playlists
 - [x] Try to parse song name and artist from filename if metadata could not be extracted
 - [x] Search not only for mp3 music
 - [ ] Support for several music folders
 - [ ] Add "likes" support
 - [ ] Add albums (from audio metadata)
 - [ ] Add smart playlists by gente (from metadata)
 - [ ] Add image (from metadata)
 - [ ] Add audio queues, so you would not need to start each song manually
 - [ ] Support for user creation (because leaking songs is pirating Arrrrr!!!) 
