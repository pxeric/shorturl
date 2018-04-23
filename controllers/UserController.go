// 用户管理类

package controllers

import (
	"shorturl/backend/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

func AddUser(c *gin.Context) {
	username, email, imcode := c.PostForm("username"), c.PostForm("email"), c.PostForm("imcode")
	allow, _ := strconv.ParseBool(c.DefaultPostForm("allow", "false"))

	if len(username) == 0 || len(email) == 0 || len(imcode) == 0 {
		ResultFail(c, models.MissParams, "缺少参数")
	} else {
		id, _ := strconv.Atoi(c.DefaultPostForm("id", "0"))
		if id <= 0 {
			user := models.User{UserName: username, Email: email, Imcode: imcode, Allow: allow}
			err := models.InsertUser(&user)
			if err == nil {
				ResultOK(c, "创建成功", "", nil)
			} else {
				ResultFail(c, models.Fail, "创建失败，原因："+err.Error())
			}
		} else {
			user := models.User{UserId: id, UserName: username, Email: email, Imcode: imcode, Allow: allow}
			err := models.UpdateUser(&user)
			if err == nil {
				ResultOK(c, "修改成功", "", nil)
			} else {
				ResultFail(c, models.NotExist, "修改失败，原因："+err.Error())
			}
		}
	}
}

func GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	list := models.GetUserList(page, size)

	ResultOK(c, "获取成功", "", list)
}

func GetUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
	if id <= 0 {
		ResultFail(c, models.MissParams, "缺少参数")
	} else {
		model := models.GetUserById(id)
		if model != nil {
			ResultOK(c, "获取成功", "", model)
		} else {
			ResultFail(c, models.NotExist, "获取失败")
		}
	}
}

func DelUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	if id <= 0 {
		ResultFail(c, models.MissParams, "缺少参数")
	} else {
		model := models.GetUserById(id)
		if model != nil {
			models.DeleteUser(id)
			ResultOK(c, "删除成功", "", nil)
		} else {
			ResultFail(c, models.Fail, "删除失败")
		}
	}
}
