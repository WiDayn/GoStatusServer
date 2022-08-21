package controller

import (
	"GoStatusServer/model"
	"GoStatusServer/response"
	"GoStatusServer/utils"
	"github.com/gin-gonic/gin"
)

func GetListDot(clients []model.Client) []model.GetListFeedback {
	var getList []model.GetListFeedback
	for _, client := range clients {
		getList = append(getList, model.GetListFeedback{
			ClientId:    client.ClientId,
			DisplayName: client.DisplayName,
			CountryCode: client.CountryCode,
		})
	}
	return getList
}

func GetList(c *gin.Context) {
	var clients []model.Client
	utils.GormDb.Find(&clients)
	response.Success(c, gin.H{"list": GetListDot(clients)}, "Success")
}
