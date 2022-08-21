package telegram

import (
	"GoStatusServer/model"
	"GoStatusServer/utils"
	"github.com/goccy/go-json"
	tele "gopkg.in/telebot.v3"
	"strconv"
	"time"
)

func Status(c tele.Context) error {
	var clients []model.Client
	utils.GormDb.Find(&clients)

	text := time.Now().Format("[2022-08-21 22:57:20]") + "\n"

	for _, client := range clients {
		text += "ID: " + client.ClientId + "\n"
		text += "显示名词: " + client.DisplayName + "\n"
		var updateRequest model.UpdateRequest
		result, _ := utils.Redisdb.Get(client.ClientId).Result()
		_ = json.Unmarshal([]byte(result), &updateRequest)

		if updateRequest.UpdateTime.Before(time.Now().Add(-time.Second * 15)) {
			text += "离线 \n"
		} else {
			text += "CPU占用: " + strconv.FormatFloat(updateRequest.DynamicInformation.CPUAvg, 'g', 3, 32) + "\n"
			text += "内存占用: " + strconv.FormatFloat(updateRequest.DynamicInformation.MemUsedPercent, 'g', 3, 32) + "\n"
			text += "上行 | 下行: " + strconv.Itoa(updateRequest.DynamicInformation.NowUpStreamDataSize) + " | " + strconv.Itoa(updateRequest.DynamicInformation.NowDownStreamDataSize) + "\n"
			text += "CT | CU | CM: " + strconv.Itoa(int(updateRequest.DynamicInformation.CT.AvgRTT)) + " | " + strconv.Itoa(int(updateRequest.DynamicInformation.CU.AvgRTT)) + " | " + strconv.Itoa(int(updateRequest.DynamicInformation.CM.AvgRTT)) + "\n"
			text += "CT 丢包率：" + strconv.FormatFloat(float64(updateRequest.DynamicInformation.CT.PacketsReceive)/float64(updateRequest.DynamicInformation.CT.PacketsSent*100), 'g', 3, 32) + "\n"
			text += "CU 丢包率：" + strconv.FormatFloat(float64(updateRequest.DynamicInformation.CU.PacketsReceive)/float64(updateRequest.DynamicInformation.CU.PacketsSent*100), 'g', 3, 32) + "\n"
			text += "CM 丢包率：" + strconv.FormatFloat(float64(updateRequest.DynamicInformation.CM.PacketsReceive)/float64(updateRequest.DynamicInformation.CM.PacketsSent*100), 'g', 3, 32) + "\n"
			text += "最后更新时间: " + updateRequest.UpdateTime.Format("2006-01-02 15:04:05") + "\n"
		}
		text += "\n"
	}

	return c.Send(text)
}
