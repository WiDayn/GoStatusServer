package main

import (
	"GoStatusServer/config"
	"GoStatusServer/controller"
	"GoStatusServer/logger"
	"GoStatusServer/model"
	"GoStatusServer/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(200)
		} else {
			ctx.Next()
		}
	}
}

func main() {
	config.Read()
	utils.SQLInit()
	utils.RedisInit()
	model.ClientInit()
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.POST("/register", controller.Register)
	r.GET("/update", controller.Update)
	r.GET("/query", controller.Query)
	r.GET("/getList", controller.GetList)
	r.GET("/getBaseInformation", controller.GetBasicInformation)
	r.GET("/queryDetail", controller.QueryDetail)
	if err := r.Run(":" + strconv.Itoa(config.Config.Port)); err != nil {
		logger.Panic("Open HTTP server error", err)
	}
}
