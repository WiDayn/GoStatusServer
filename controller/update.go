package controller

import (
	"GoStatusServer/logger"
	"GoStatusServer/model"
	"GoStatusServer/utils"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

//设置websocket
//CheckOrigin防止跨站点的请求伪造
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Update(c *gin.Context) {
	lastTimeUpdatePing := time.Now()
	lastTimeClearPing := time.Now()
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
	}(ws)
	for {
		var updateRequest model.UpdateRequest
		err := ws.ReadJSON(&updateRequest)
		if err != nil {
			logger.Error("Read updateRequest json error", nil)
			return
		}

		if lastTimeClearPing.Before(time.Now().Add(-time.Hour * 24)) {
			// 清空redis，重新记录
			marshal, _ := json.Marshal(model.PingRecords{
				CT: model.PingRecord{PacketsReceive: 0, PacketsSent: 0},
				CU: model.PingRecord{PacketsReceive: 0, PacketsSent: 0},
				CM: model.PingRecord{PacketsReceive: 0, PacketsSent: 0},
			})
			utils.Redisdb.Set(updateRequest.ClientId+"/PingRecords", marshal, time.Hour*20480)
			lastTimeClearPing = time.Now()
		}

		if lastTimeUpdatePing.Before(time.Now().Add(-time.Minute * 5)) {
			//每五分钟记录一次
			if result, err := utils.Redisdb.Get(updateRequest.ClientId + "/PingRecords").Result(); err == nil {
				var pingRecords model.PingRecords
				err = json.Unmarshal([]byte(result), &pingRecords)
				if err != nil {
					return
				}
				pingRecords.CT = model.PingRecord{
					PacketsSent:    pingRecords.CT.PacketsSent + updateRequest.DynamicInformation.CT.PacketsSent,
					PacketsReceive: pingRecords.CT.PacketsReceive + updateRequest.DynamicInformation.CT.PacketsReceive,
					Time:           append(pingRecords.CT.Time, updateRequest.UpdateTime.Format("2006-01-02 15:04:05")),
					AvgRTT:         append(pingRecords.CT.AvgRTT, updateRequest.DynamicInformation.CT.AvgRTT),
				}
				pingRecords.CU = model.PingRecord{
					PacketsSent:    pingRecords.CU.PacketsSent + updateRequest.DynamicInformation.CU.PacketsSent,
					PacketsReceive: pingRecords.CU.PacketsReceive + updateRequest.DynamicInformation.CU.PacketsReceive,
					Time:           append(pingRecords.CU.Time, updateRequest.UpdateTime.Format("2006-01-02 15:04:05")),
					AvgRTT:         append(pingRecords.CU.AvgRTT, updateRequest.DynamicInformation.CU.AvgRTT),
				}
				pingRecords.CM = model.PingRecord{
					PacketsSent:    pingRecords.CM.PacketsSent + updateRequest.DynamicInformation.CM.PacketsSent,
					PacketsReceive: pingRecords.CM.PacketsReceive + updateRequest.DynamicInformation.CM.PacketsReceive,
					Time:           append(pingRecords.CM.Time, updateRequest.UpdateTime.Format("2006-01-02 15:04:05")),
					AvgRTT:         append(pingRecords.CM.AvgRTT, updateRequest.DynamicInformation.CM.AvgRTT),
				}
				marshal, err := json.Marshal(pingRecords)
				if err != nil {
					return
				}
				utils.Redisdb.Set(updateRequest.ClientId+"/PingRecords", marshal, time.Hour*20480)
				lastTimeUpdatePing = time.Now()
			}

			if result, err := utils.Redisdb.Get(updateRequest.ClientId + "/BasicRecords").Result(); err == nil {
				var basicRecords model.BasicRecords
				err = json.Unmarshal([]byte(result), &basicRecords)
				if err != nil {
					return
				}

				basicRecords.MemUsedPercent = append(basicRecords.MemUsedPercent, updateRequest.DynamicInformation.MemUsedPercent)
				basicRecords.CPUAvg = append(basicRecords.CPUAvg, updateRequest.DynamicInformation.CPUAvg)
				basicRecords.DiskPercent = append(basicRecords.DiskPercent, updateRequest.DynamicInformation.DiskInformation.Percent)
				basicRecords.NowUpStreamDataSize = append(basicRecords.NowUpStreamDataSize, updateRequest.DynamicInformation.NowUpStreamDataSize)
				basicRecords.NowDownStreamDataSize = append(basicRecords.NowDownStreamDataSize, updateRequest.DynamicInformation.NowDownStreamDataSize)
				basicRecords.Time = append(basicRecords.Time, updateRequest.UpdateTime.Format("2006-01-02 15:04:05"))

				marshal, err := json.Marshal(basicRecords)
				if err != nil {
					return
				}
				utils.Redisdb.Set(updateRequest.ClientId+"/BasicRecords", marshal, time.Hour*20480)
				lastTimeUpdatePing = time.Now()
			}
		}
		js, _ := json.Marshal(updateRequest)
		utils.Redisdb.Set(updateRequest.ClientId, js, time.Hour*20480)
	}
}
