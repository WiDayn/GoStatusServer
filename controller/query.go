package controller

import (
	"GoStatusServer/logger"
	"GoStatusServer/utils"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
)

type QueryRequest struct {
	ClientsId []string
}

type QueryFeedback struct {
	ClientId                string
	CPUAvg                  float64
	MemAll                  uint64
	MenFree                 uint64
	MenUsed                 uint64
	MemUsedPercent          float64
	TotalDownStreamDataSize uint64
	TotalUpStreamDataSize   uint64
	NowDownStreamDataSize   int
	NowUpStreamDataSize     int
	DiskTotal               uint64
	DiskUsed                uint64
	DiskPercent             uint64
}

func QueryFeedbackDto(client UpdateRequest) QueryFeedback {
	return QueryFeedback{
		ClientId:                client.ClientId,
		TotalDownStreamDataSize: client.DynamicInformation.TotalDownStreamDataSize,
		TotalUpStreamDataSize:   client.DynamicInformation.TotalUpStreamDataSize,
		NowUpStreamDataSize:     client.DynamicInformation.NowUpStreamDataSize,
		NowDownStreamDataSize:   client.DynamicInformation.NowDownStreamDataSize,
		CPUAvg:                  client.DynamicInformation.CPUAvg,
		MemAll:                  client.DynamicInformation.MemAll,
		MenFree:                 client.DynamicInformation.MenFree,
		MenUsed:                 client.DynamicInformation.MenUsed,
		MemUsedPercent:          client.DynamicInformation.MemUsedPercent,
		DiskTotal:               client.DynamicInformation.DiskInformation.Total,
		DiskUsed:                client.DynamicInformation.DiskInformation.Used,
		DiskPercent:             client.DynamicInformation.DiskInformation.Percent,
	}
}

func Query(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("Upgrade websocket error", err)
	}
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			logger.Warning("Close websocket error", err)
			panic(err)
		}
	}(ws) //返回前关闭
	var queryRequest QueryRequest
	if err := ws.ReadJSON(&queryRequest); err != nil {
		logger.Error("Read updateRequest json error", nil)
		return
	}
	for {
		var queryFeedback []QueryFeedback
		for _, clientId := range queryRequest.ClientsId {
			res, _ := utils.Redisdb.Keys(clientId).Result()
			var client UpdateRequest
			err := json.Unmarshal([]byte(res[0]), &client)
			if err != nil {
				logger.Error("Read redis json error", nil)
				return
			}
			queryFeedback = append(queryFeedback, QueryFeedbackDto(client))
		}
		err := ws.WriteJSON(queryFeedback)
		if err != nil {
			logger.Error("Write updateRequest json error", nil)
			return
		}
	}
}
