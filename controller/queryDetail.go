package controller

import (
	"GoStatusServer/model"
	"GoStatusServer/response"
	"GoStatusServer/utils"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
)

func QueryDetail(c *gin.Context) {
	ClientId := c.Query("ClientId")

	var client model.Client
	var pingRecords model.PingRecords
	var basicRecords model.BasicRecords

	utils.GormDb.Model(&model.Client{}).Where("client_id = ?", ClientId).First(&client)
	result, _ := utils.Redisdb.Get(ClientId + "/PingRecords").Result()
	if err := json.Unmarshal([]byte(result), &pingRecords); err != nil {
		return
	}

	result, _ = utils.Redisdb.Get(ClientId + "/BasicRecords").Result()
	if err := json.Unmarshal([]byte(result), &basicRecords); err != nil {
		return
	}

	response.Success(c, gin.H{"client": client, "pingRecords": pingRecords, "basicRecords": basicRecords}, "Success")
}
