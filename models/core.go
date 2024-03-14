package model

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var store redis.Store
var dsn string

const (
	MaxDBOpenConn   = 1800
	MaxDBIdleConn   = 10
	UserSession     = "usersession"
	CommonTextMax   = 500
	CommonStringMax = 255   // the limit set when creating table
	MySQLTextMax    = 16383 // Text max is 65535bytes, 16383 characters using utf8mb4
	UUIDLen         = 36
	PwdMD5Len       = 32
	SaltLen         = 64
	ThreadsPerPage  = 10
	PostsPerPage    = 10
	CommentsPerPage = 10
)

func init() {

	cfg, err := ini.Load("./config/app.ini")
	if err != nil {
		panic("config not load ")
	}

	dsn = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Section("mysql").Key("user").String(),
		cfg.Section("mysql").Key("password").String(),
		cfg.Section("mysql").Key("ip").String(),
		cfg.Section("mysql").Key("port").String(),
		cfg.Section("mysql").Key("database").String())
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		QueryFields:            true,
		DisableAutomaticPing:   true,
		SkipDefaultTransaction: true, //禁用GORM在执行CREATE/UPDATE/DELETE时自动开启事务
	})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get sql.DB")
	}
	sqlDB.SetMaxIdleConns(MaxDBIdleConn)
	sqlDB.SetMaxOpenConns(MaxDBOpenConn)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)

	// cookie based session, store in redis, the cookie contains the session id but no data
	kp1 := cfg.Section("redis").Key("keypair1").String()
	kp2 := cfg.Section("redis").Key("keypair2").String()
	maxidle, _ := cfg.Section("redis").Key("maxidle").Int()
	store, _ = redis.NewStore(
		maxidle,
		cfg.Section("redis").Key("network").String(),
		cfg.Section("redis").Key("ip").String(),
		cfg.Section("redis").Key("password").String(),
		[]byte(kp1),
		[]byte(kp2),
	)

	// cookie config also applies to session data in redis
	maxage, _ := cfg.Section("cookie").Key("maxage").Int()
	secure, _ := cfg.Section("cookie").Key("secure").Bool()
	httponly, _ := cfg.Section("cookie").Key("httponly").Bool()
	samesite, _ := cfg.Section("cookie").Key("samesite").Int()
	store.Options(sessions.Options{
		Path:     cfg.Section("cookie").Key("path").String(),
		Domain:   cfg.Section("cookie").Key("domain").String(),
		MaxAge:   maxage,
		Secure:   secure,
		HttpOnly: httponly,
		SameSite: http.SameSite(samesite),
	})

}

// func GetDB() *gorm.DB {
// 	return db
// }

func GetSessionFunc(sessionNames []string) gin.HandlerFunc {
	return sessions.SessionsMany(sessionNames, store)
}
func GetSession(c *gin.Context) sessions.Session {
	return sessions.DefaultMany(c, UserSession)
}
func SleepConnsKiller(db *gorm.DB) {
	var sleepIds []int
	db.Raw("SELECT id FROM INFORMATION_SCHEMA.processlist").Scan(&sleepIds)
	for id := range sleepIds {
		db.Exec("KILL ?", id)
	}
}
