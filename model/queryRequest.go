package model

type QueryRequest struct {
	ClientsId   string `json:"ClientId"`
	DisplayName string `json:"DisplayName"`
	CountryCode string `json:"CountryCode"`
}

type QueryFeedback struct {
	ClientId                string
	DisplayName             string
	CountryCode             string
	CPUAvg                  float64
	MemAll                  string
	MenFree                 string
	MenUsed                 string
	MemUsedPercent          float64
	TotalDownStreamDataSize string
	TotalUpStreamDataSize   string
	NowDownStreamDataSize   string
	NowUpStreamDataSize     string
	DiskTotal               string
	DiskUsed                string
	DiskPercent             uint64
	Online                  bool
	CU                      PingStatus
	CT                      PingStatus
	CM                      PingStatus
}
