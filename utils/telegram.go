package utils

import (
	"GoStatusServer/config"
	"GoStatusServer/logger"
	"time"

	tele "gopkg.in/telebot.v3"
)

var TelegramBot *tele.Bot

func InitTelegramBot() {
	if !config.Config.TelegramBot.Enable {
		return
	}
	var err error

	pref := tele.Settings{
		Token:  config.Config.TelegramBot.BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	TelegramBot, err = tele.NewBot(pref)
	if err != nil {
		logger.Error("Error to start telegram bot", err)
		return
	}
}

func SendTelegramNotify(text string) {
	if !config.Config.TelegramBot.Enable {
		return
	}

	_, err := TelegramBot.Send(&tele.User{ID: config.Config.TelegramBot.NotifyID}, time.Now().Format("[2006-01-02 15:04:05]  ")+text)
	if err != nil {
		logger.Error("Error to send message", err)
	}
}
