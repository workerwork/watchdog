package task

import (
	"time"
)

//系统信息上报redis
func Statsystem() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		<-ticker.C
		writeRedis("RPUSH", "eGWstat_check", "system")
		readRedis("BRPOP", "eGWstat_readflag", "0")
	}
}

func Stater() {
	/*检测系统*/
	go Statsystem()
}
