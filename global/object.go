package global

type UsableObject interface {
	Usable() bool
	Key() string
}

const (
	KeyMysql           = "mysql"
	KeyConnPool        = "connPool"
	KeyLogger          = "logger"
	KeyDispatcher      = "dispatcher"
	KeyPlayListService = "playListService"
)
