package controller

import (
	"GoStatusServer/model"
	"GoStatusServer/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	ClientId         string
	DisplayName      string
	BasicInformation BasicInformation
}

type BasicInformation struct {
	IP             string
	CPUs           []CPU
	CPUPhysicalCnt int
	CPULogicalCnt  int
	OS             string
	Hostname       string
	Status         string  `json:"status"`
	Country        string  `json:"country"`
	CountryCode    string  `json:"countryCode"`
	Region         string  `json:"region"`
	RegionName     string  `json:"regionName"`
	City           string  `json:"city"`
	Zip            string  `json:"zip"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	Timezone       string  `json:"timezone"`
	Isp            string  `json:"isp"`
	Org            string  `json:"org"`
	As             string  `json:"as"`
	Query          string  `json:"query"`
}

type CPU struct {
	CPUName  string
	CPUCores int32
}

func Register(c *gin.Context) {
	var registerRequest RegisterRequest
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
}
