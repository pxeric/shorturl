package controllers

import (
	"shorturl/backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ResultFail(c *gin.Context, code models.ReturnCode, message string) {
	c.JSON(http.StatusOK, gin.H{"code": code, "message": message})
	c.Abort()
}

func ResultOK(c *gin.Context, message, shorturl string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": models.Success, "message": message, "shorturl": shorturl, "data": data})
	c.Abort()
}
