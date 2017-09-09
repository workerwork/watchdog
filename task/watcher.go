package task

import (
	"time"
)

//监控进程
func Watchprocess(checkcmd, checker, executecmd string) {
	ticker := time.NewTicker(2 * time.Second)

	for {
		<-ticker.C
		//time := <-ticker.C
		//layout := "2006-01-02 15:04:05"
		//t := time.Format(layout)
		//fmt.Println("watchprocess:", t)
		check, _ := Exec(checkcmd)
		if check[0] != checker {
			Run(executecmd)
		}
		//return
	}
}

func Watcher() {
	/*egw_manage进程监控*/
	go Watchprocess("ps -ef|awk '/egw_manage$/ {print $8}'", "/root/eGW/OMC/egw_manage", "spawn-fcgi -a 127.0.0.1 -p 4001 -f /root/eGW/OMC/egw_manage")
	/*egw_report进程监控*/
	go Watchprocess("ps -ef|awk '/egw_report$/ {print $8}'", "/root/eGW/OMC/egw_report", "/root/eGW/OMC/egw_report &")
}
