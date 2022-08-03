package controller

import (
	"GoStatusServer/logger"
	"GoStatusServer/utils"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type UpdateRequest struct {
	ClientId           string
	DynamicInformation DynamicInformation
}

type DynamicInformation struct {
	CPUAvg                  float64
	MemAll                  uint64
	MenFree                 uint64
	MenUsed                 uint64
	MemUsedPercent          float64
	TotalDownStreamDataSize uint64
	TotalUpStreamDataSize   uint64
	NowDownStreamDataSize   int
	NowUpStreamDataSize     int
	DiskInformation         DiskInformation
}

type DiskInformation struct {
	Total   uint64
	Used    uint64
	Percent uint64
}

//设置websocket
//CheckOrigin防止跨站点的请求伪造
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Update(c *gin.Context) {
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
	for {
		var updateRequest UpdateRequest
		err := ws.ReadJSON(&updateRequest)
		if err != nil {
			logger.Error("Read updateRequest json error", nil)
			return
		}
		js, _ := json.Marshal(updateRequest)

		utils.Redisdb.Set(updateRequest.ClientId, js, time.Minute*10)
		println(updateRequest.DynamicInformation.CPUAvg)
	}
}
