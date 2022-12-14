package watcher

import (
	"GoStatusServer/config"
	"GoStatusServer/logger"
	"GoStatusServer/model"
	"GoStatusServer/utils"
	"github.com/goccy/go-json"
	"strconv"
	"time"
)

type OnlineStatusWatcher struct {
	ClientList []clientStatus
}

type clientStatus struct {
	ClientId      string
	Status        bool // true -> online
	HighCPUCount  int
	HighMemCount  int
	HighDiskCount int
	Count         int
}

var DefaultOnlineStatusWatcher OnlineStatusWatcher

func (onlineStatusWatcher *OnlineStatusWatcher) Add(clientId string) {
	onlineStatusWatcher.ClientList = append(onlineStatusWatcher.ClientList, clientStatus{ClientId: clientId, Status: false})
}

func (onlineStatusWatcher *OnlineStatusWatcher) Run() {
	logger.Info("监控线程启动", nil)
	go func() {
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
				if updateRequest.DynamicInformation.CPUAvg > config.Config.Watcher.CPUPercent {
					onlineStatusWatcher.ClientList[key].HighCPUCount++
				}
				if updateRequest.DynamicInformation.MemUsedPercent > config.Config.Watcher.MemPercent {
					onlineStatusWatcher.ClientList[key].HighMemCount++
				}
				if updateRequest.DynamicInformation.DiskInformation.Percent > config.Config.Watcher.DiskPercent {
					onlineStatusWatcher.ClientList[key].HighDiskCount++
				}
				if clientStatus.HighCPUCount > 120 {
					utils.SendTelegramNotify(updateRequest.DisplayName + " 近期CPU负载较高!\n 当前CPU占用率为:" + strconv.FormatFloat(updateRequest.DynamicInformation.CPUAvg, 'g', 3, 32))
					clearCount(&onlineStatusWatcher.ClientList[key])
				}
				if clientStatus.HighMemCount > 120 {
					utils.SendTelegramNotify(updateRequest.DisplayName + " 近期内存负载较高!\n 当前内存占用率为:" + strconv.FormatFloat(updateRequest.DynamicInformation.MemUsedPercent, 'g', 3, 32))
					clearCount(&onlineStatusWatcher.ClientList[key])
				}
				if clientStatus.HighDiskCount > 120 {
					utils.SendTelegramNotify(updateRequest.DisplayName + " 近期硬盘占用较高!\n 当前硬盘占用率为:" + strconv.Itoa(int(updateRequest.DynamicInformation.DiskInformation.Percent)))
					clearCount(&onlineStatusWatcher.ClientList[key])
				}

				if clientStatus.Count > 320 {
					clearCount(&onlineStatusWatcher.ClientList[key])
				}
				clientStatus.Count++
			}
			time.Sleep(time.Second * 3)
		}
	}()
}

func clearCount(clientStatus *clientStatus) {
	clientStatus.Count = 0
	clientStatus.HighCPUCount = 0
	clientStatus.HighMemCount = 0
	clientStatus.HighDiskCount = 0
}
