package models

import (
	"shorturl/backend/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/core"
	"github.com/xormplus/xorm"
	"net/url"
)

var db *xorm.Engine

func InitDB() {
	username := utils.AppConfig.DBUser
	password := utils.AppConfig.DBPwd
	hostname := utils.AppConfig.DBHost
	dbName := utils.AppConfig.DBName
	protocol := "tcp"
	timezone := "Asia/Shanghai"
	charset := "utf8mb4"
	collation := "utf8mb4_bin"

	dsn := username + ":" + password +
		"@" + protocol + "(" + hostname + ")" +
		"/" + dbName +
		"?parseTime=True&loc=" + url.PathEscape(timezone) +
		"&charset=" + charset + "&collation=" + collation

	var err error
	db, err = xorm.NewEngine("mysql", dsn)
	if err != nil {
		panic(err)
	}

	//db.SetMaxIdleConns(utils.AppConfig.MaxIdle) //设置闲置连接数
	//db.SetMaxOpenConns(utils.AppConfig.MaxOpen) //设置最大打开的连接数，默认为0表示不限制

	if utils.AppConfig.DebugLog {
		db.ShowSQL(true)
		db.Logger().SetLevel(core.LOG_DEBUG)
	} else {
		db.Logger().SetLevel(core.LOG_INFO)
	}

	err = db.Sync2(new(Urlmap), new(User), new(Statistics))
	if err != nil {
		panic(err)
	}
}
