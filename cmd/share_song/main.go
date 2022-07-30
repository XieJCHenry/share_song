package main

import (
	"fmt"
	"share_song/music/music_library"
	"share_song/music/play_list"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	Main2()
}

func Main1() {
	config, err := loadConfig()
	if err != nil {
		panic(fmt.Sprintf("load config failed, err=%s", err.Error()))
	}

	//ctx := context.Background()
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		return
	}
	logger := zapLog.Sugar()

	e := gin.Default()

	mysqlDsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.MySql.UserName,
		config.MySql.Password,
		config.MySql.Host,
		config.MySql.Port,
		config.MySql.DataBase,
	)
	mysqlDb, err := gorm.Open(mysql.Open(mysqlDsn))
	if err != nil {
		panic(fmt.Sprintf("connect mysql failed, err=%s", err))
	}

	//mongoDsn := fmt.Sprintf("mongodb://%s:%d", config.Mongo.Host, config.Mongo.Port)
	//mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDsn))
	//if err != nil {
	//	panic(fmt.Sprintf("connect mongo failed, err=%s", err))
	//}

	//userService := user.NewService(logger, mysqlDb, mongoClient)
	//userController := user.NewController(logger, userService)
	//userController.Register(e)

	musicLibraryService := music_library.NewMusicLibraryService(logger, mysqlDb)
	musicLibraryController := music_library.NewController(logger, musicLibraryService)
	musicLibraryController.Register(e)

	playerListService := play_list.NewService(logger, nil, musicLibraryService)
	playListController := play_list.NewController(logger, playerListService)
	playListController.Register(e)

	e.Run(":8874")
}
