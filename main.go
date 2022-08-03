package main

import (
	"GoStatusServer/config"
	"GoStatusServer/controller"
	"GoStatusServer/logger"
	"GoStatusServer/model"
	"GoStatusServer/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

//设置websocket
//CheckOrigin防止跨站点的请求伪造
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//websocket实现
func ping(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			panic(err)
		}
	}(ws) //返回前关闭
	for {
		//读取ws中的数据
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		//写入ws数据
		err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

func main() {
	config.Read()
	utils.SQLInit()
	utils.RedisInit()
	model.ClientInit()
	r := gin.Default()
	r.GET("/ping", ping)
	r.POST("/register", controller.Register)
	r.GET("/update", controller.Update)
	r.GET("/query", controller.Query)
	r.GET("/getList", controller.GetList)
	if err := r.Run(":" + strconv.Itoa(config.Config.Port)); err != nil {
		logger.Panic("Open HTTP server error", err)
	}
}
