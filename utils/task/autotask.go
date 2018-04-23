package task

import (
	"shorturl/backend/caches"
	"shorturl/backend/models"
	"fmt"
	"strings"
	"time"
)

func Start() {
	//开启一个定时任务
	fmt.Println("定时任务已开启")
	ticker := time.NewTicker(time.Second * 60)
	go func() {
		for _ = range ticker.C {
			keys, err := caches.RedisCli.GetKeys("total_*") //如：total_pc_mgymav
			if err == nil {
				for _, key := range keys {
					arr := strings.Split(key, "_")
					platform := arr[1]
					shorturl := arr[2]
					total, err := caches.RedisCli.GetInt(key)
					if err == nil {
						err := models.AddTotal(shorturl, platform, total)
						if err == nil {
							caches.RedisCli.Del(key)
						}
					}
				}
			}
		}
	}()
}
