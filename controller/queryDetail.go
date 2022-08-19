package controller

import (
	"GoStatusServer/model"
	"GoStatusServer/response"
	"GoStatusServer/utils"
	"github.com/gin-gonic/gin"
)

func QueryDetail(c *gin.Context) {
	ClientId := c.Query("ClientId")

	var client model.Client

	utils.GormDb.Model(&model.Client{}).Where("client_id = ?", ClientId).First(&client)

	response.Success(c, gin.H{"client": client}, "Success")
}
