package user

type UserDto struct {
	// mysql fields
	InstanceId string `json:"instance_id,omitempty"`
	Name       string `json:"name,omitempty"`
	Phone      string `json:"phone,omitempty"`
}

type User struct {
	// mysql fields
	InstanceId string `json:"instance_id,omitempty"`
	Name       string `json:"name,omitempty"`
	Phone      string `json:"phone,omitempty"`

	//
	Status Status `json:"status,omitempty"`
	// 本次登录对歌单操作过的歌曲
	OperatedSongs   []*SongOperation `json:"operatedSongs,omitempty" bson:"operated_songs"`
	addedSongsMap   map[string]int
	deletedSongsMap map[string]int
}

func (u *User) UpdateOperatedSongs(songs []string, op SongOpr) {

	if op == Add {
		for _, id := range songs {
			if _, ok := u.addedSongsMap[id]; !ok {
				u.OperatedSongs = append(u.OperatedSongs, &SongOperation{
					SongId:    id,
					Operation: op,
				})
				u.addedSongsMap[id] = len(u.OperatedSongs) - 1
			}
		}
	} else if op == Delete {
		for _, id := range songs {
			if _, ok := u.deletedSongsMap[id]; !ok {
				u.OperatedSongs = append(u.OperatedSongs, &SongOperation{
					SongId:    id,
					Operation: op,
				})
				u.deletedSongsMap[id] = len(u.OperatedSongs) - 1
			}
		}
	}
}
