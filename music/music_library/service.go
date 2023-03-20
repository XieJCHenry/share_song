package music_library

import "share_song/music"

type MusicLibraryService interface {
	BatchAddSong(songs []*music.Song) (*BatchSongResult, error)
	BatchRemoveSong(songIds []string) (*BatchSongResult, error)
	//BatchUpdateSong(songs []*music.Song) (*BatchSongResult, error)
	SearchSongs(fields []string, page, pageSize int) ([]music.Song, error)
	SearchSongByCondition(fields []string, query map[string]interface{}, page, pageSize int) ([]music.Song, error)
}
