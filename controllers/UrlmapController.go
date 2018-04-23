package controllers

import (
	"shorturl/backend/caches"
	"shorturl/backend/models"
	"shorturl/backend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func EncodeUrlmap(c *gin.Context) {
	longurl, iosurl, androidurl, customurl := c.DefaultPostForm("longurl", c.Query("longurl")), c.DefaultPostForm("iosurl", c.Query("iosurl")), c.DefaultPostForm("androidurl", c.Query("androidurl")), c.DefaultPostForm("customurl", c.Query("customurl"))

	if len(longurl) == 0 {
		ResultFail(c, models.MissParams, "longurl参数错误")
	} else {
		if !(utils.StringContains(longurl, "http://", "https://")) {
			longurl = "http://" + longurl
		}
		if len(iosurl) > 0 && !(utils.StringContains(iosurl, "http://", "https://")) {
			iosurl = "http://" + iosurl
		}
		if len(androidurl) > 0 && !(utils.StringContains(androidurl, "http://", "https://")) {
			androidurl = "http://" + androidurl
		}

		//根据长链获取短链
		model := models.GetUrlmapByLongUrl(longurl)

		shorturl := ""
		if len(customurl) > 0 && (model == nil || customurl != model.CustomUrl) {
			tmp := models.GetUrlmapByCustomUrl(customurl)
			if tmp != nil {
				ResultFail(c, models.Exist, "http://hostname/"+customurl+"已存在，请重新输入")
				return
			} else {
				shorturl = customurl
			}
		}

		//判断短链是否已存在，存在则直接返回
		if model == nil {
			if shorturl == "" { //随机生成短链
				for {
					tempshorturl := utils.GenerateShortUrl(6)
					tmp := models.GetUrlmapByCustomUrl(tempshorturl)
					if tmp == nil {
						shorturl = tempshorturl
						break
					}
				}
			}

			//写入缓存和数据库
			urlmap := models.Urlmap{ShortUrl: shorturl, LongUrl: longurl, IOSUrl: iosurl, AndroidUrl: androidurl, CustomUrl: customurl}
			err := models.InsertUrlmap(&urlmap)
			if err == nil {
				models.InsertStatistics(&models.Statistics{Url: shorturl})
				caches.RedisCli.Set(shorturl, urlmap, -1)
				ResultOK(c, "生成成功", "http://hostname/"+shorturl, nil)
			} else {
				ResultFail(c, models.Fail, "生成短链失败，原因："+err.Error())
			}
		} else {
			//ios链接和android链接不相同时修改数据库中的值
			if iosurl != model.IOSUrl || androidurl != model.AndroidUrl || customurl != model.CustomUrl {
				model.IOSUrl = iosurl
				model.AndroidUrl = androidurl
				model.CustomUrl = customurl
				models.UpdateUrlmap(model)
			}
			//写入缓存
			caches.RedisCli.Set(model.ShortUrl, model, -1)
			//写一份自定义链接的缓存
			if len(model.CustomUrl) > 0 {
				caches.RedisCli.Set(model.CustomUrl, model, -1)
			}
			ResultOK(c, "生成成功", "http://hostname/"+model.ShortUrl, nil)
		}
	}
}

func DecodeUrlmap(c *gin.Context) {
	shorturl := c.DefaultPostForm("shorturl", c.Query("shorturl"))
	shorturl = utils.StringReplace(shorturl, "", "http://hostname/")
	if len(shorturl) == 0 {
		ResultFail(c, models.MissParams, "缺少参数")
	} else {
		model := models.GetUrlmapByCustomUrl(shorturl)
		if model != nil {
			ResultOK(c, "获取成功", "", model)
		} else {
			ResultFail(c, models.NotExist, "记录不存在")
		}
	}
}

func RedirectUrl(c *gin.Context) {
	shorturl := c.Param("shorturl")
	if len(shorturl) == 0 {
		c.Redirect(http.StatusFound, "http://www.baidu.com")
	} else {
		var model *models.Urlmap
		err := caches.RedisCli.GetObject(shorturl, &model)
		if err != nil {
			model = models.GetUrlmapByCustomUrl(shorturl)
			caches.RedisCli.Set(shorturl, model, -1)
		}

		if model == nil {
			c.Redirect(http.StatusFound, "http://www.baidu.com")
		} else {
			url := model.LongUrl
			//判断什么平台访问的
			platform := utils.GetPlatform(c.GetHeader("User-Agent"))
			if platform == "ios" {
				url = model.IOSUrl
			} else if platform == "android" {
				url = model.AndroidUrl
			}

			//创建一个线程记录跳转次数
			go func(surl, plat string) {
				caches.RedisCli.Incr("total_" + plat + "_" + shorturl) //将跳转次数写入redis，然后通过定时任务再写入mysql
			}(shorturl, platform)

			//302跳转
			c.Redirect(http.StatusFound, url)
		}
	}
}

func GetUrlmapList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	shorturl := utils.StringReplace(c.Query("shorturl"), "", "http://hostname/")

	list := models.GetUrlmapList(page, size, shorturl)

	ResultOK(c, "获取成功", "", list)
}

func GetUrlmap(c *gin.Context) {
	shorturl := c.Query("shorturl")
	if shorturl == "" {
		ResultFail(c, models.MissParams, "缺少参数")
	} else {
		model := models.GetUrlmapByShortUrl(shorturl)
		if model != nil {
			ResultOK(c, "获取成功", "", model)
		} else {
			ResultFail(c, models.NotExist, "获取失败")
		}
	}
}
