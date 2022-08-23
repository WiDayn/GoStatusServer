package controller

import (
	"GoStatusServer/config"
	"GoStatusServer/model"
	"GoStatusServer/response"
	"GoStatusServer/utils"
	"GoStatusServer/watcher"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

func Register(c *gin.Context) {
	if c.Query("SecretKey") != config.Config.SecretKey {
		response.Fail(c, nil, "SecretKey Wrong")
		return
	}

	var registerRequest model.RegisterRequest
	if err := c.BindJSON(&registerRequest); err != nil {
		return
	}
	var client model.Client
	client.ClientId = registerRequest.ClientId
	client.DisplayName = registerRequest.DisplayName
	for _, cpu := range registerRequest.BasicInformation.CPUs {
		client.CPU = append(client.CPU, model.CPU{
			ClientId: registerRequest.ClientId,
			CPUName:  cpu.CPUName,
			CPUCores: cpu.CPUCores,
		})
	}
	client.IP = registerRequest.BasicInformation.IP
	client.OS = registerRequest.BasicInformation.OS
	client.CPULogicalCnt = registerRequest.BasicInformation.CPULogicalCnt
	client.CPUPhysicalCnt = registerRequest.BasicInformation.CPUPhysicalCnt
	client.Hostname = registerRequest.BasicInformation.Hostname
	client.Country = registerRequest.BasicInformation.Country
	client.As = registerRequest.BasicInformation.As
	client.CountryCode = registerRequest.BasicInformation.CountryCode
	client.Region = registerRequest.BasicInformation.Region
	client.RegionName = registerRequest.BasicInformation.RegionName
	client.City = registerRequest.BasicInformation.City
	client.Zip = registerRequest.BasicInformation.Zip
	client.Lat = registerRequest.BasicInformation.Lat
	client.Lon = registerRequest.BasicInformation.Lon
	client.Timezone = registerRequest.BasicInformation.Timezone
	client.Isp = registerRequest.BasicInformation.Isp
	client.Org = registerRequest.BasicInformation.Org

	utils.GormDb.Session(&gorm.Session{FullSaveAssociations: true}).FirstOrCreate(&model.Client{ClientId: registerRequest.ClientId}).Updates(&client)
	utils.Redisdb.Set(client.ClientId+"/PingRecords", "{}", time.Hour*20480)
	utils.Redisdb.Set(client.ClientId+"/BasicRecords", "{}", time.Hour*20480)
	watcher.DefaultOnlineStatusWatcher.Add(client.ClientId)
}
