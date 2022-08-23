package telegram

import (
	"GoStatusServer/config"
	"GoStatusServer/utils"
	tele "gopkg.in/telebot.v3"
)

func RegisterController() {
	if !config.Config.TelegramBot.Enable {
		return
	}

	utils.TelegramBot.Handle("/ping", func(c tele.Context) error {
		return c.Send("pong")
	})

	utils.TelegramBot.Handle("/start", func(c tele.Context) error {
		return c.Send("服务器启动！ 响应正常")
	})

	utils.TelegramBot.Handle("/status", Status)

	go utils.TelegramBot.Start()

	utils.SendTelegramNotify("监控服务器已启动!")
}
