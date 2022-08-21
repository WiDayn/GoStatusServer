package watcher

import (
	"GoStatusServer/config"
	"GoStatusServer/logger"
	"GoStatusServer/model"
	"GoStatusServer/utils"
	"github.com/goccy/go-json"
	"time"
)

type OnlineStatusWatcher struct {
	ClientList []clientStatus
}

type clientStatus struct {
	ClientId string
	Status   bool // true -> online
}

var DefaultOnlineStatusWatcher OnlineStatusWatcher

func (onlineStatusWatcher *OnlineStatusWatcher) Add(clientId string) {
	onlineStatusWatcher.ClientList = append(onlineStatusWatcher.ClientList, clientStatus{ClientId: clientId, Status: false})
}

func (onlineStatusWatcher *OnlineStatusWatcher) Run() {
	logger.Info("监控线程启动", nil)
	for {
		for key, clientStatus := range onlineStatusWatcher.ClientList {
			var updateRequest model.UpdateRequest
			var nowStatus bool
			result, _ := utils.Redisdb.Get(clientStatus.ClientId).Result()
			_ = json.Unmarshal([]byte(result), &updateRequest)

			if updateRequest.UpdateTime.Before(time.Now().Add(-time.Second * 15)) {
				nowStatus = false
			} else {
				nowStatus = true
			}

			if clientStatus.Status != nowStatus {
				if config.Config.TelegramBot.Enable {
					if nowStatus {
						utils.SendTelegramNotify(updateRequest.DisplayName + " 已上线")
					} else {
						utils.SendTelegramNotify(updateRequest.DisplayName + " 已离线")
					}
				}
				onlineStatusWatcher.ClientList[key].Status = nowStatus
			}
		}
		time.Sleep(time.Second * 2)
	}
}
