package controller

import (
	"GoStatusServer/config"
	"GoStatusServer/model"
	"GoStatusServer/response"
	"github.com/gin-gonic/gin"
)

func GetNotifyList(c *gin.Context) {
	var notifyList model.NotifyList

	notifyList.Telegram = config.Config.TelegramBot.Enable

	response.Success(c, gin.H{"NotifyList": notifyList}, "Success")
}
