package main

import (
	"fmt"
	"time"
	"watchdog.git/task"
)

func main() {
	fmt.Println("\n", time.Now(), "watchdog进程启动:")
	fmt.Println("***********************************************************")
	task.Getcpuall()
	task.Watcher()
	task.Manager()
	task.Stater()
	task.Checker()
	/*主线程循环*/
	for {
		time.Sleep(60 * time.Second)
	}
	fmt.Println("OK")
}
