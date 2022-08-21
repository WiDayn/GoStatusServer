package model

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
