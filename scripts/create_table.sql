CREATE TABLE `t_user`
(
    `instance_id` char(11)    NOT NULL COMMENT '用户ID',
    `name`        varchar(16) NOT NULL DEFAULT '' COMMENT '账户名',
    `phone`       char(13)    NOT NULL DEFAULT '' COMMENT '密码',
    PRIMARY KEY (`instance_id`)
) ENGINE=InnoDB AUTO_INCREMENT=136 DEFAULT CHARSET=utf8 COMMENT='用户';



CREATE TABLE `t_song`
(
    `instance_id` char(11)     NOT NULL COMMENT '歌曲ID',
    `name`        varchar(15)  NOT NULL DEFAULT '' COMMENT '歌曲名',
    `artists`     varchar(511) NOT NULL DEFAULT '' COMMENT '演唱者',
    `length`      int          NOT NULL DEFAULT 0 COMMENT '时长，单位秒',
    `lyric`       varchar(255) NOT NULL DEFAULT '' COMMENT '歌词文件路径',
    `source`      varchar(255) NOT NULL DEFAULT '' COMMENT '本地文件或url',
    `cover`       varchar(256) NOT NULL DEFAULT '' COMMENT '封面地址',
    `from`        varchar(32)  NOT NULL DEFAULT '' COMMENT '来源类型',
    `binary_size` bigint       NOT NULL DEFAULT 0 COMMENT '二进制大小，单位byte',
    `delete_flag` tinyint(1) NOT NULL DEFAULT 0 COMMENT '删除标记',
    PRIMARY KEY (`instance_id`)
) ENGINE=InnoDB AUTO_INCREMENT=136 DEFAULT CHARSET=utf8 COMMENT='歌曲信息';