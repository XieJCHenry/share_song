package music_library

import (
	"fmt"
	"share_song/music"
	"share_song/utils/uuid"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	songTableName = "t_song"
)

type service struct {
	logger *zap.SugaredLogger
	mysql  *gorm.DB
}

func NewMusicLibraryService(logger *zap.SugaredLogger, mysql *gorm.DB) MusicLibraryService {
	logger.Errorf("曲库目前仍无法通过接口调用添加歌曲")
	return &service{
		logger: logger,
		mysql:  mysql,
	}
}

func (s *service) BatchAddSong(songs []*music.Song) (*BatchSongResult, error) {
	var names []string
	filterMap := make(map[string]int)
	for i, song := range songs {
		if _, ok := filterMap[song.Name]; !ok {
			filterMap[song.Name] = i
			names = append(names, song.Name)
		}
	}

	if len(names) <= 0 {
		return &BatchSongResult{}, nil
	}

	var searchResultList []music.Song
	result := s.mysql.Table(songTableName).
		Select("instance_id", "name", "delete_flag").
		Where("name in ?", names).
		Find(&searchResultList)
	if result.Error != nil {
		s.logger.Errorf("search songs failed, err=%s", result.Error)
		return nil, fmt.Errorf("数据库查询歌曲错误：%s", result.Error)
	}
	var batchAddList []*music.Song
	var batchUpdateList []string
	failData := make(map[string]interface{})
	if len(searchResultList) > 0 {
		resultMap := make(map[string]int)
		for i, song := range searchResultList {
			resultMap[song.Name] = i
		}
		for songName, index := range filterMap {
			idx, isExists := resultMap[songName]
			searchResult := searchResultList[idx]
			if !isExists {
				newSong := songs[index]
				newSong.InstanceId = uuid.GenerateWithLength(uuid.InstanceIdLength)
				batchAddList = append(batchAddList, newSong)
			} else if searchResult.DeleteFlag == true {
				batchUpdateList = append(batchUpdateList, searchResult.InstanceId)
			} else {
				failData[searchResult.Name] = fmt.Sprintf("《%s》已经存在", searchResult.Name)
			}
		}
	} else {
		for i := range songs {
			newSong := songs[i]
			newSong.InstanceId = uuid.GenerateWithLength(uuid.InstanceIdLength)
			batchAddList = append(batchAddList, newSong)
		}
	}

	if len(batchAddList) > 0 {
		insertList := convertToDtoList(batchAddList)
		// add songs
		result = s.mysql.Table(songTableName).CreateInBatches(insertList, 20)
		if result.Error != nil {
			s.logger.Errorf("create songs failed, err=%s", result.Error)
			return nil, fmt.Errorf("数据库新建歌曲失败：%s", result.Error)
		}
	}

	if len(batchUpdateList) > 0 {
		// update delete_flag
		result = s.mysql.Table(songTableName).
			Where("instance_id in ?", batchUpdateList).
			Update("delete_flag", false)
		if result.Error != nil {
			s.logger.Errorf("update delete_mark failed, err=%s", result.Error)
			return nil, fmt.Errorf("未知错误：%s", result.Error)
		}
	}

	return &BatchSongResult{
		InsertCount: len(batchAddList) + len(batchUpdateList),
		FailData:    failData,
	}, nil
}

func (s *service) BatchRemoveSong(songIds []string) (*BatchSongResult, error) {
	if len(songIds) <= 0 {
		return &BatchSongResult{}, nil
	}

	result := s.mysql.Table(songTableName).
		Where("instance_id in ?", songIds).
		Update("delete_flag", true)
	if result.Error != nil {
		s.logger.Errorf("set delete_mark=true failed, err=%s", result.Error)
		return nil, fmt.Errorf("删除歌曲错误：%s", result.Error)
	}

	return &BatchSongResult{
		RemoveCount: len(songIds),
	}, nil
}

func (s *service) BatchUpdateSong(songs []*music.Song) (*BatchSongResult, error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) SearchSongs(fields []string, page, pageSize int) ([]music.Song, error) {
	if page < 0 {
		page = 1
	}
	if pageSize < 0 {
		pageSize = 50
	}
	offset := (page - 1) * pageSize
	var songs []music.Song
	result := s.mysql.Table(songTableName).Select(fields).
		Where("delete_flag = false").
		Offset(offset).Limit(pageSize).
		Find(&songs)
	if result.Error != nil {
		s.logger.Errorf("search songs failed, err=%s", result.Error)
		return nil, fmt.Errorf("查找歌曲失败，错误：%s", result.Error)
	}

	return songs, nil
}

func (s *service) SearchSongByCondition(fields []string, query map[string]interface{}, page, pageSize int) ([]music.Song, error) {
	if !validateSearchQuery(query) {
		return nil, fmt.Errorf("查询参数query不合法，当前仅支持instanceId查找或为空")
	}
	var songIds []string
	if val, ok := query["instance_id"].(string); ok {
		songIds = append(songIds, val)
	} else if val, ok := query["instance_id"].([]string); ok {
		songIds = val
	}

	offset := (page - 1) * pageSize
	var songs []music.Song
	result := s.mysql.Table(songTableName).Select(fields).
		Where("instance_id in ?", songIds).
		Offset(offset).Limit(pageSize).
		Find(&songs)
	if result.Error != nil {
		s.logger.Errorf("search songs failed, err=%s", result.Error)
		return nil, fmt.Errorf("查找歌曲失败，错误：%s", result.Error)
	}
	return songs, nil
}

func validateSearchQuery(query map[string]interface{}) bool {
	if len(query) <= 0 {
		return true
	}
	if _, ok := query["instance_id"]; ok && len(query) == 1 {
		return true
	}
	return false
}

func encodeAuthors(authors []string) string {
	return strings.Join(authors, ",")
}

func decodeAuthors(authors string) []string {
	return strings.Split(authors, ",")
}

func convertToDtoList(songs []*music.Song) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(songs))
	for _, song := range songs {
		result = append(result, map[string]interface{}{
			"instance_id": song.InstanceId,
			"name":        song.Name,
			"artists":     encodeAuthors(song.Artists),
			"length":      song.Length,
			"lyric":       song.Lyric,
			"source":      song.Source,
			"cover":       song.Cover,
			"delete_flag": song.DeleteFlag,
			"from":        song.From,
			"binary_size": song.BinarySize,
		})
	}

	return result
}
