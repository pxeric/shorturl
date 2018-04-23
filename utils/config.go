package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type (
	Configuration struct {
		IsDebug    bool
		HttpPort   string
		UploadPath string
		DBHost     string
		DBUser     string
		DBPwd      string
		DBName     string
		RedisHost  string
		RedisPwd   string
		RedisDB    int
		DebugLog   bool
	}
)

var AppConfig Configuration

func InitConfig() {
	loadAppConfig()
}

func loadAppConfig() {
	//for _, val := range os.Environ() {
	//	fmt.Println(val)
	//}

	//获取环境变量
	environment := os.Getenv("GOLANG_ENVIRONMENT")
	confPath := "conf/conf.json"
	//若设置了环境变量则加载对应的配置文件
	if environment != "" {
		confPath = "conf/conf." + environment + ".json"
	}

	data, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Fatalf("[loadConfig]: %s\n", err)
	}
	err = json.Unmarshal(data, &AppConfig)
	if err != nil {
		log.Fatalf("[loadAppConfig]: %s\n", err)
	}
}
