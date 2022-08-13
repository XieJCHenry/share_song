package user

import (
	"context"
	"fmt"
	"regexp"
	"share_song/internal/wbsocket"
	"share_song/utils/uuid"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	MySqlTableName = "t_user"

	MongoDbName  = "share_song"
	MongoColName = "user_song_operation"

	phoneRegex = "^1(3\\d|4[5-9]|5[0-35-9]|6[2567]|7[0-8]|8\\d|9[0-35-9])\\d{8}$"
)

type service struct {
	logger            *zap.SugaredLogger
	loginNotifier     chan string
	onlinePool        *sync.Map
	timeOutOnlinePool *sync.Map
	mysql             *gorm.DB
	mongoClient       *mongo.Client

	connPool *wbsocket.Pool
}

func NewService(logger *zap.SugaredLogger, db *gorm.DB, mongoClient *mongo.Client) Service {
	return &service{
		logger:            logger,
		loginNotifier:     make(chan string),
		onlinePool:        &sync.Map{},
		timeOutOnlinePool: &sync.Map{},
		mysql:             db,
		mongoClient:       mongoClient,
	}
}

func (s *service) Login(ctx context.Context, userName string, phone string) (*User, error) {
	user := &User{}

	tbl := s.mysql.Table(MySqlTableName).Select("instance_id", "name", "phone")
	if len(userName) > 0 {
		tbl.Where("name = ?", userName)
	}
	if len(phone) > 0 {
		tbl.Or("phone = ?", phone)
	}
	result := tbl.First(user)
	err := result.Error
	if err != nil {
		s.logger.Errorf("login check user exists failed, err=%s", err)
		return nil, fmt.Errorf("数据库查询用户错误：%s", err.Error())
	}

	// already online
	if onlineUser, ok := s.onlinePool.Load(user.InstanceId); ok {
		return onlineUser.(*User), nil
	}

	cursor, err := s.mongoClient.Database(MongoDbName).Collection(MongoColName).Find(ctx, bson.D{
		{"user_id", user.InstanceId},
	})
	if err != nil {
		s.logger.Errorf("search user operated songs failed, err=%s", err.Error())
		return nil, fmt.Errorf("数据库查询用户歌曲操作列表失败：%s", err.Error())
	}
	for cursor.Next(ctx) {
		var tempUser User
		err := cursor.Decode(&tempUser)
		if err != nil {
			return nil, fmt.Errorf("解析用户歌曲操作列表失败：%s", err.Error())
		}
		user.OperatedSongs = append(user.OperatedSongs, tempUser.OperatedSongs...)
	}
	defer cursor.Close(ctx)

	user.Status = Online
	s.onlinePool.Store(user.InstanceId, user)
	s.loginNotifier <- user.InstanceId
	s.logger.Infof("user %s(%s) logined", user.Name, user.InstanceId)

	return user, nil
}

func (s *service) Logout(ctx context.Context, instanceId string, userName string) error {
	_, err := s.loadUser(instanceId, "", "")
	if err != nil {
		s.logger.Errorf("logout check user exists failed, err=%s", err)
		return err
	}

	if online, ok := s.onlinePool.Load(instanceId); ok {
		onlineUser := online.(*User)
		onlineUser.Status = TimeOutOffline
		s.onlinePool.Delete(onlineUser.InstanceId)
		s.timeOutOnlinePool.Store(onlineUser.InstanceId, onlineUser)
		s.logger.Infof("user %s(%s) timeout logout", onlineUser.Name, onlineUser.InstanceId)

		// timeout offline. 如果在1分内登录回来，则恢复登录状态
		go func(logoutGuyId string) {
			timer := time.AfterFunc(60*time.Second, func() {
				s.logger.Infof("user %s logout", logoutGuyId)
				s.trulyLogOut(ctx, logoutGuyId)
			})

			select {
			case instanceId := <-s.loginNotifier:
				if instanceId == logoutGuyId {
					if timeOutOfflineGuy, loaded := s.timeOutOnlinePool.LoadAndDelete(logoutGuyId); loaded {
						timeOutGuy := timeOutOfflineGuy.(*User)
						timeOutGuy.Status = Online
						s.onlinePool.Store(logoutGuyId, timeOutGuy)
						timer.Stop()
						return
					}
				}
			}
		}(onlineUser.InstanceId)

	}
	return nil
}

func (s *service) trulyLogOut(ctx context.Context, userId string) {
	timeOutOnlineGuy, ok := s.timeOutOnlinePool.Load(userId)
	if ok {
		defer func() {
			s.timeOutOnlinePool.Delete(userId)
		}()

		timeOutGuy := timeOutOnlineGuy.(*User)
		// 将玩家歌单操作记录到mongo todo 需要测试歌曲加载和持久化的正确性
		updateResult, err := s.mongoClient.Database(MongoDbName).Collection(MongoColName).UpdateOne(ctx, bson.D{
			{"user_id", timeOutGuy.InstanceId},
		}, bson.D{
			{"user_id", timeOutGuy.InstanceId},
			{"operated_songs", timeOutGuy.OperatedSongs},
		})
		if err != nil {
			s.logger.Errorf("update user %s song operations failed, err=%s", timeOutGuy.InstanceId, err.Error())
			return
		}
		s.logger.Infof("update user %s song operations: matched %d upsert %d modified %d",
			timeOutGuy.InstanceId, updateResult.MatchedCount, updateResult.UpsertedCount, updateResult.ModifiedCount)
	}
}

func (s *service) RegisterAccount(ctx context.Context, userName string, phone string) (*User, error) {
	// check exists
	tbl := s.mysql.Table(MySqlTableName)
	selectFields := []string{"instance_id"}
	if len(userName) > 0 {
		tbl.Where("name = ?", userName)
		selectFields = append(selectFields, "name")
	}
	if len(phone) > 0 {
		tbl.Or("phone = ?", phone)
		selectFields = append(selectFields, "phone")
	}

	var selectedUsers []UserDto
	searchResult := s.mysql.Table(MySqlTableName).Select(selectFields).Find(selectedUsers)
	err := searchResult.Error
	if err != nil {
		s.logger.Errorf("check user is existed failed, err=%s", err)
		return nil, fmt.Errorf("数据库查询错误：%s", err.Error())
	}
	for _, user := range selectedUsers {
		if user.Name == userName {
			return nil, fmt.Errorf("用户名（%s）已被使用", userName)
		}
		if user.Phone == phone {
			return nil, fmt.Errorf("手机号（%s）已被使用", phone)
		}
	}

	// create new
	newUser := &UserDto{
		InstanceId: uuid.GenerateWithLength(uuid.InstanceIdLength),
		Name:       userName,
		Phone:      phone,
	}
	result := s.mysql.Table(MySqlTableName).Create(newUser)
	if result.Error != nil {
		s.logger.Errorf("create user(%s:%s) failed, err=%s", userName, phone, result.Error)
		return nil, fmt.Errorf("创建用户错误，err=%s", result.Error.Error())
	}
	s.logger.Infof("new user register %s(%s)", userName, newUser.InstanceId)
	return &User{
		InstanceId: uuid.GenerateWithLength(uuid.InstanceIdLength),
		Name:       userName,
		Phone:      phone,
	}, nil
}

func (s *service) CancelAccount(ctx context.Context, instanceId, userName, phone string) error {
	// 删除db中账户信息
	existsUser, err := s.loadUser(instanceId, userName, phone)
	if err != nil {
		s.logger.Errorf("check user %s(%s) is existed failed, err=%s", userName, instanceId, err)
		return fmt.Errorf("加载用户%s失败，err=%s", userName, err.Error())
	}

	// 删除在线列表或超时列表中的信息
	s.onlinePool.Delete(instanceId)
	s.timeOutOnlinePool.Delete(instanceId)

	if len(existsUser.InstanceId) > 0 {
		result := s.mysql.Table(MySqlTableName).Delete(existsUser)
		if result.Error != nil {
			s.logger.Errorf("delete user(%s:%s) failed, err=%s", userName, phone, result.Error)
			return result.Error
		}

		deleteResult, err := s.mongoClient.Database(MongoDbName).Collection(MongoColName).DeleteOne(ctx, bson.D{
			{"user_id", instanceId},
		})
		if err != nil {
			s.logger.Errorf("mongo delte user %s(%s) song operations failed, err=%s", instanceId, userName, err.Error())
			return err
		}
		s.logger.Infof("mongo delete user %s(%s) song operations, count %d", instanceId, userName, deleteResult.DeletedCount)
	}

	s.logger.Infof("user %s(%s) canceled", userName, instanceId)

	return nil
}

func (s *service) loadUser(instanceId, userName, phone string) (*User, error) {
	limited := &User{
		InstanceId: instanceId,
		Name:       userName,
		Phone:      phone,
	}
	result := &User{}
	dbResult := s.mysql.Table(MySqlTableName).Where(limited).First(result)
	if dbResult.Error != nil {
		s.logger.Errorf("check user is existed failed, err=%s", dbResult.Error)
		return nil, dbResult.Error
	}
	return result, nil
}

func (s *service) SearchOnlineUsers(ctx context.Context, query map[string]interface{}, withTimeOut bool) ([]*User, error) {

	if !validateSearchQuery(query) {
		return nil, fmt.Errorf("查询参数query不合法，当前仅支持instanceId查找或为空")
	}

	if user := s.searchOnlineByInstanceId(query, withTimeOut); user != nil {
		return []*User{user}, nil
	}

	var users []*User
	s.onlinePool.Range(func(key, value any) bool {
		users = append(users, value.(*User))
		return true
	})
	if withTimeOut {
		s.timeOutOnlinePool.Range(func(key, value any) bool {
			users = append(users, value.(*User))
			return true
		})
	}
	return users, nil
}

func (s *service) searchOnlineByInstanceId(query map[string]interface{}, withTimeOut bool) *User {
	if instanceId, ok := query["instanceId"]; ok {
		if onlineUser, loaded := s.onlinePool.Load(instanceId.(string)); loaded {
			return onlineUser.(*User)
		}

		if withTimeOut {
			if timeOutUser, loaded := s.timeOutOnlinePool.Load(instanceId.(string)); loaded {
				return timeOutUser.(*User)
			}
		}
	}
	return nil
}

func validateSearchQuery(query map[string]interface{}) bool {
	if len(query) <= 0 {
		return true
	}
	if _, ok := query["instanceId"]; ok && len(query) == 1 {
		return true
	}
	return false
}

func (s *service) GetOnlineToken(ctx context.Context, loginKey string) (string, error) {

	tbl := s.mysql.Table(MySqlTableName).Select([]string{"instance_id", "name", "phone"})

	isPhone, err := checkIsPhone(loginKey)
	if err != nil {
		return "", fmt.Errorf("请输入正确的用户名或手机号")
	}
	if isPhone {
		tbl.Where("phone = ?", loginKey)
	} else {
		tbl.Where("name = ?", loginKey)

	}

	user := &UserDto{}
	result := tbl.First(user)
	if result.Error != nil {
		s.logger.Errorf("login check user exists failed, err=%s", result.Error)
		return "", fmt.Errorf("数据库查询用户错误：%s", result.Error.Error())
	}
	return user.InstanceId, nil
}

func checkIsPhone(key string) (bool, error) {
	pattern, err := regexp.Compile(phoneRegex)
	if err != nil {
		return false, err
	}
	if len(key) == 13 {
		if pattern.MatchString(key) {
			return true, nil
		}
	}
	return false, nil
}
