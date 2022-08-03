package controller

import (
	"GoStatusServer/model"
	"GoStatusServer/response"
	"GoStatusServer/utils"
	"github.com/gin-gonic/gin"
)

type GetListFeedback struct {
	ClientId    string
	DisplayName string
	CountryCode string `gorm:"varchar(40)"`
}

func GetListDot(clients []model.Client) []GetListFeedback {
	var getList []GetListFeedback
	for _, client := range clients {
		getList = append(getList, GetListFeedback{
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
