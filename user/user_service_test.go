package user

import (
	"fmt"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestConnectMysql(t *testing.T) {
	host := "127.0.0.1"
	port := 3306
	userName := "root"
	password := "xjc241003"
	dbName := "share_song"
	timeout := 10

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%d", userName, password, host, port, dbName, timeout)

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		t.Fatalf("open mysql failed, err=%s", err.Error())
		return
	}
	_ = db
}
