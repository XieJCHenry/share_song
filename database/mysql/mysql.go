package mysql

import (
	"share_song/global"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Client struct {
	logger *zap.SugaredLogger
	mysql  *gorm.DB
}

func NewClient(logger *zap.SugaredLogger, mysql *gorm.DB) *Client {
	return &Client{
		logger: logger,
		mysql:  mysql,
	}
}

func (c *Client) Usable() bool {
	return c.mysql != nil
}

func (c *Client) Key() string {
	return global.KeyMysql
}
