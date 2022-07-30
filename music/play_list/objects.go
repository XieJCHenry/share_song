package play_list

import "share_song/music"

type AddSongsRequest struct {
	UserId  string   `json:"userId,omitempty"`
	SongIds []string `json:"songIds,omitempty"`
}

type AddSongsResponse struct {
	Songs []music.Song `json:"songs,omitempty"`
}

type RemoveSongsRequest struct {
	UserId  string   `json:"userId,omitempty"`
	SongIds []string `json:"songIds,omitempty"`
}

type RemoveSongsResponse struct {
	Songs []music.Song `json:"songs,omitempty"`
}

type GetCurrentSongsResponse struct {
	Songs []music.Song `json:"songs,omitempty"`
}

type UpgradeUserRequest struct {
	OnlineToken string `json:"onlineToken"`
}

type UpgradeUserResponse struct {
	OnlineToken string `json:"onlineToken"`
	Msg         string `json:"msg"`
}
