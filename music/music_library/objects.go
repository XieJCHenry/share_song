package music_library

import "share_song/music"

type OperateType string

const (
	And OperateType = "and"
	Or  OperateType = "or"
)

type SearchSongCondition struct {
	Conditions map[string]interface{}
	Operator   OperateType
}

type BatchSongResult struct {
	InsertCount int
	UpdateCount int
	RemoveCount int
	FailData    map[string]interface{}
}

type BatchAddSongRequest struct {
	Songs []*music.Song `json:"songs,omitempty"`
}

type BatchAddSongResponse struct {
	Result *BatchSongResult `json:"result,omitempty"`
}

type BatchRemoveSongRequest struct {
	SongIds []string `json:"songIds,omitempty"`
}

type BatchRemoveSongResponse struct {
	Result *BatchSongResult `json:"result,omitempty"`
}

type SearchSongsRequest struct {
	Fields   []string `json:"fields,omitempty"`
	Page     int      `json:"page,omitempty"`
	PageSize int      `json:"pageSize,omitempty"`
}

type SearchSongsResponse struct {
	Songs    []music.Song `json:"songs,omitempty"`
	Page     int          `json:"page,omitempty"`
	PageSize int          `json:"pageSize,omitempty"`
}
