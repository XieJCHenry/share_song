package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	mysql2 "share_song/database/mysql"
	"share_song/global"
	"share_song/hello"
	"share_song/internal/wbsocket"
	logger2 "share_song/logger"
	"share_song/music/music_library"
	"share_song/music/play_list"
	"share_song/music/play_list_v2"
	"share_song/protocol"
	"share_song/user_v2"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		panic(fmt.Sprintf("load config failed, err=%s", err.Error()))
	}

	global.Init()

	e := gin.Default()

	ctx := context.Background()
	_ = ctx
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		return
	}
	logger := logger2.New(zapLog.Sugar())

	routeDispatcher := protocol.NewDispatcher()
	hello.Register(routeDispatcher)
	play_list_v2.Register(routeDispatcher)

	mysqlDsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.MySql.UserName,
		config.MySql.Password,
		config.MySql.Host,
		config.MySql.Port,
		config.MySql.DataBase,
	)
	mysqlDb, err := gorm.Open(mysql.Open(mysqlDsn))
	_ = mysqlDb
	if err != nil {
		panic(fmt.Sprintf("connect mysql failed, err=%s", err))
	}
	mysqlClient := mysql2.NewClient(logger.Sugared(), mysqlDb)

	connPool := wbsocket.NewConnectionPool(logger.Sugared())

	global.SetGlobalObject(logger)
	global.SetGlobalObject(routeDispatcher)
	global.SetGlobalObject(mysqlClient)
	global.SetGlobalObject(connPool)

	userServiceV2 := user_v2.NewUserServiceV2(logger.Sugared())
	userController := user_v2.NewController(userServiceV2)
	userController.Register(e)

	musicLibraryService := music_library.NewMusicLibraryService(logger.Sugared(), mysqlDb)
	musicLibraryController := music_library.NewController(logger.Sugared(), musicLibraryService)
	musicLibraryController.Register(e)

	playListService := play_list.NewService(logger.Sugared(), nil, musicLibraryService)
	playListController := play_list.NewController(logger.Sugared(), playListService)
	playListController.Register(e)

	global.SetGlobalObject(playListService)

	e.Run(":8874")
}