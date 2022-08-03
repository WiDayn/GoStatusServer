package model

import (
	"GoStatusServer/utils"
	"time"
)

type Client struct {
	ClientId                string `gorm:"type:varchar(40); primaryKey"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DisplayName             string `gorm:"varchar(20)"`
	IP                      string `gorm:"varchar(20)"`
	CPU                     []CPU  `gorm:"foreignKey:ClientId"`
	CPUPhysicalCnt          int    `gorm:"int(6)"`
	CPULogicalCnt           int    `gorm:"int(6)"`
	OS                      string `gorm:"varchar(120)"`
	Hostname                string `gorm:"varchar(120)"`
	CPUavg                  float64
	MemAll                  uint64
	MenFree                 uint64
	MenUsed                 uint64
	MemUsedPercent          float64
	TotalDownStreamDataSize uint64
	TotalUpStreamDataSize   uint64
	NowDownStreamDataSize   int
	NowUpStreamDataSize     int
	Country                 string  `gorm:"varchar(40)"`
	CountryCode             string  `gorm:"varchar(40)"`
	Region                  string  `gorm:"varchar(40)"`
	RegionName              string  `gorm:"varchar(40)"`
	City                    string  `gorm:"varchar(40)"`
	Zip                     string  `gorm:"varchar(40)"`
	Lat                     float64 `gorm:"varchar(40)"`
	Lon                     float64 `gorm:"varchar(40)"`
	Timezone                string  `gorm:"varchar(40)"`
	Isp                     string  `gorm:"varchar(40)"`
	Org                     string  `gorm:"varchar(40)"`
	As                      string  `gorm:"varchar(40)"`
	Query                   string  `gorm:"varchar(40)"`
}

type CPU struct {
	ClientId string `gorm:"varchar(40)"`
	CPUName  string `gorm:"varchar(120)"`
	CPUCores int32  `gorm:"int(6)"`
}

func ClientInit() {
	if err := utils.GormDb.AutoMigrate(&Client{}); err != nil {
		panic(err)
	}
	if err := utils.GormDb.AutoMigrate(&CPU{}); err != nil {
		panic(err)
	}
}
