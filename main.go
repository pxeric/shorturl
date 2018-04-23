package main

import (
	"shorturl/backend/caches"
	"shorturl/backend/controllers"
	"shorturl/backend/models"
	"shorturl/backend/utils"
	"shorturl/backend/utils/task"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func main() {
	//加载配置文件
	utils.InitConfig()

	// 初始化router
	router := gin.Default()

	// 初始化数据库
	models.InitDB()

	//初始化redis
	caches.Init()

	//开启定时任务
	task.Start()

	//session存redis
	store, _ := sessions.NewRedisStoreWithDB(10, "tcp", utils.AppConfig.RedisHost, utils.AppConfig.RedisPwd, strconv.Itoa(utils.AppConfig.RedisDB), []byte("secret"))

	//设置session过期时间
	store.Options(sessions.Options{
		MaxAge:   1800, //单位为秒，设置30分钟
		Path:     "/",
		HttpOnly: true,
	})

	//session
	router.Use(sessions.Sessions("urlmap", store))

	// 用户登录
	router.POST("/api/user/login", controllers.Login)
	router.POST("/api/user/getimcode", controllers.GetIMCode)
	router.POST("/api/user/logout", controllers.Logout)

	//短链相关接口
	router.GET("/api/redirect/:shorturl", controllers.RedirectUrl)
	router.GET("/urlmap/encode", controllers.EncodeUrlmap)
	router.POST("/urlmap/encode", controllers.EncodeUrlmap)
	router.GET("/urlmap/decode", controllers.DecodeUrlmap)
	router.POST("/urlmap/decode", controllers.DecodeUrlmap)

	admin := router.Group("/api/admin")
	admin.Use(LoginFilter())
	{
		// 用户管理
		admin.GET("getuser", controllers.GetUser)
		admin.GET("userlist", controllers.GetUserList)
		admin.POST("adduser", controllers.AddUser)
		admin.POST("deluser", controllers.DelUser)
		admin.GET("urlmaplist", controllers.GetUrlmapList)
		admin.GET("urlmap", controllers.GetUrlmap)
		admin.POST("urlmap/encode", controllers.EncodeUrlmap)
	}

	//启动服务
	router.Run(utils.AppConfig.HttpPort)
}

func LoginFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sessions.Default(c)
		user := sess.Get("userid")
		if user != nil && user != "" {
			c.Next()
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": models.NoLogin, "message": "未登录"})
		c.Abort()
	}
}
