package controller

import (
	"GoStatusServer/logger"
	"GoStatusServer/model"
	"GoStatusServer/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"strconv"
	"time"
)

func QueryFeedbackDto(client model.UpdateRequest) model.QueryFeedback {
	return model.QueryFeedback{
		ClientId:                client.ClientId,
		DisplayName:             client.DisplayName,
		CountryCode:             client.CountryCode,
		TotalDownStreamDataSize: strconv.Itoa(int(client.DynamicInformation.TotalDownStreamDataSize/1024/1024/1024)) + " GB",
		TotalUpStreamDataSize:   strconv.Itoa(int(client.DynamicInformation.TotalUpStreamDataSize/1024/1024/1024)) + " GB",
		NowUpStreamDataSize:     strconv.Itoa(client.DynamicInformation.NowUpStreamDataSize) + " Kbp/s",
		NowDownStreamDataSize:   strconv.Itoa(client.DynamicInformation.NowDownStreamDataSize) + " Kbp/s",
		CPUAvg:                  ParseFloat(client.DynamicInformation.CPUAvg),
		MemAll:                  strconv.Itoa(int(client.DynamicInformation.MemAll/1024/1024)) + " MB",
		MenFree:                 strconv.Itoa(int(client.DynamicInformation.MenFree/1024/1024)) + " MB",
		MenUsed:                 strconv.Itoa(int(client.DynamicInformation.MenUsed/1024/1024)) + " MB",
		MemUsedPercent:          ParseFloat(client.DynamicInformation.MemUsedPercent),
		DiskTotal:               strconv.Itoa(int(client.DynamicInformation.DiskInformation.Total)) + " GB",
		DiskUsed:                strconv.Itoa(int(client.DynamicInformation.DiskInformation.Used)) + " GB",
		DiskPercent:             client.DynamicInformation.DiskInformation.Percent,
		Online:                  client.Online,
		CT:                      client.DynamicInformation.CT,
		CU:                      client.DynamicInformation.CU,
		CM:                      client.DynamicInformation.CM,
	}
}

func ParseFloat(num float64) float64 {
	num, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", num), 64)
	return num
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
			return
		}
	}(ws) //返回前关闭
	var queryRequest []model.QueryRequest
	err = ws.ReadJSON(&queryRequest)
	for err != nil {
		// logger.Error("Read updateRequest json error", nil)
		err = ws.ReadJSON(&queryRequest)
	}
	for {
		var queryFeedback []model.QueryFeedback
		for _, query := range queryRequest {
			if res, err := utils.Redisdb.Get(query.ClientsId).Result(); err != redis.Nil {
				var client model.UpdateRequest
				err := json.Unmarshal([]byte(res), &client)
				if err != nil {
					logger.Error("Read redis json error", nil)
					return
				}
				if client.UpdateTime.After(time.Now().Add(-time.Second * 15)) {
					client.Online = true
				} else {
					client.Online = false
				}
				queryFeedback = append(queryFeedback, QueryFeedbackDto(client))
			}

		}
		err := ws.WriteJSON(queryFeedback)
		if err != nil {
			logger.Error("Write updateRequest json error", nil)
			err := ws.Close()
			if err != nil {
				logger.Warning("Close websocket error", err)
				return
			}
			return
		}
		time.Sleep(time.Second * 1)
	}
}
