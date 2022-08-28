package config

import (
	"fmt"
	"github.com/go-ini/ini"
	uuid "github.com/satori/go.uuid"
	"log"
	"os"
)

type config struct {
	Debug          bool
	SecretKey      string  `comment:"不填会随机生成"`
	RecordInterval float64 `comment:"记录的间隔(粒度)，单位分钟"`
	Port           int
	SQLServer      SQLServer
	RedisServer    RedisServer
	TelegramBot    TelegramBot
	Watcher        Watcher
}

type SQLServer struct {
	IP       string
	Port     int
	Username string
	Password string
	Db       string
	Charset  string
	Loc      string
}

type RedisServer struct {
	Address  string
	Password string `comment:"未设密码可直接不填"`
	DB       int
}

type TelegramBot struct {
	Enable   bool
	BotToken string `comment:"TelegramBot的Token，不需要带前缀bot'"`
	NotifyID int64  `comment:"支持user, group, chanel, 以逗号隔开"`
}

type Watcher struct {
	Enable      bool
	CPUPercent  float64 `comment:"0-100，支持小数点， 超过100即为不提醒"`
	MemPercent  float64 `comment:"0-100，支持小数点， 超过100即为不提醒"`
	DiskPercent uint64  `comment:"0-100，不支持小数点， 超过100即为不提醒"`
}

var Config config

// Read
// this method will read the config.ini in the same directory and load it into the Config struct
// you should call this method before you use config.Config, or it will be the default value.
func Read() {
	if !isExist("./config.ini") {
		if _, err := os.Create("./config.ini"); err != nil {
			log.Println("[PANIC]Can not create config.ini")
			log.Fatalf(err.Error())
		}
	}

	if cfg, err := ini.Load("./config.ini"); err != nil {
		log.Println("[PANIC] Can not load config.ini")
		log.Fatalf(err.Error())
	} else {
		Config = config{
			Debug:          false,
			SecretKey:      uuid.NewV4().String(),
			RecordInterval: 5,
			Port:           12345,
			SQLServer: SQLServer{
				IP:       "127.0.0.1",
				Port:     3306,
				Username: "root",
				Password: "TqOMaR6zrMQUm3BG",
				Db:       "gostatus",
				Charset:  "utf8mb4",
				Loc:      "Local",
			},
			RedisServer: RedisServer{
				Address:  "127.0.0.1:6379",
				Password: "",
				DB:       0,
			},
			Watcher: Watcher{
				Enable:      false,
				CPUPercent:  80,
				MemPercent:  80,
				DiskPercent: 80,
			},
		}

		if err = ini.MapTo(&Config, "./config.ini"); err != nil {
			log.Println("[PANIC] config.ini error")
			log.Fatalf(err.Error())
		}

		if err := ini.ReflectFrom(cfg, &Config); err != nil {
			log.Println("[PANIC] config.ini error")
			log.Fatalf(err.Error())
		}

		if err := cfg.SaveTo("./config.ini"); err != nil {
			log.Println("[PANIC] Can not sava config.ini")
			log.Fatalf(err.Error())
		}
	}
	log.Println("[SUCCESS] Success load config.ini")
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}
