package controllers

import (
	"shorturl/backend/models"
	"shorturl/backend/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// 登录
func Login(c *gin.Context) {
	//需完成登录逻辑
	ResultFail(c, models.Fail, "登录失败")
}

// 退出
func Logout(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Clear()
	sess.Save()
	ResultOK(c, "退出成功", "", nil)
}
