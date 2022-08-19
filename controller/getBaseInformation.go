package controller

import (
	"GoStatusServer/model"
	"GoStatusServer/response"
	"GoStatusServer/utils"
	"github.com/gin-gonic/gin"
)

func GetBasicInformation(c *gin.Context) {
	ClientId, _ := c.Params.Get("ClientId")

	var Client model.Client

	utils.GormDb.Model(&model.Client{ClientId: ClientId}).First(&Client)

	response.Success(c, gin.H{"Client": Client}, "Success")
}
