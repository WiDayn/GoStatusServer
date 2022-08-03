package controller

import (
	"GoStatusServer/model"
	"GoStatusServer/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	ClientId    string
	DisplayName string
	BasicInf    BasicInformation
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
	for _, cpu := range registerRequest.BasicInf.CPUs {
		client.CPU = append(client.CPU, model.CPU{
			ClientId: registerRequest.ClientId,
			CPUName:  cpu.CPUName,
			CPUCores: cpu.CPUCores,
		})
	}
	client.IP = registerRequest.BasicInf.IP
	client.OS = registerRequest.BasicInf.OS
	client.CPULogicalCnt = registerRequest.BasicInf.CPULogicalCnt
	client.CPUPhysicalCnt = registerRequest.BasicInf.CPUPhysicalCnt
	client.Hostname = registerRequest.BasicInf.Hostname
	client.Country = registerRequest.BasicInf.Country
	client.As = registerRequest.BasicInf.As
	client.CountryCode = registerRequest.BasicInf.CountryCode
	client.Region = registerRequest.BasicInf.Region
	client.RegionName = registerRequest.BasicInf.RegionName
	client.City = registerRequest.BasicInf.City
	client.Zip = registerRequest.BasicInf.Zip
	client.Lat = registerRequest.BasicInf.Lat
	client.Lon = registerRequest.BasicInf.Lon
	client.Timezone = registerRequest.BasicInf.Timezone
	client.Isp = registerRequest.BasicInf.Isp
	client.Org = registerRequest.BasicInf.Org

	utils.GormDb.Session(&gorm.Session{FullSaveAssociations: true}).FirstOrCreate(&model.Client{ClientId: registerRequest.ClientId}).Updates(&client)
}
