package task

import (
	"strings"
	"time"
)

//话单管理
func Managecdr() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		time := <-ticker.C
		//layout := "2006-01-02 15:04:05"
		layout := "15:04"
		t := time.Format(layout)
		//fmt.Println("managecdr:", t)
		if t == "00:00" {
			Run("cd /root/eGW/CDR/cdrDat; ls -lt |awk '{if(NR>=10){print $9}}' |xargs rm -rf")
		}
		cdrdata, _ := Exec("cd /root/eGW/CDR; ls -lt *.dat |awk '{if(NR>=2){print $9}}'")
		cdrlen := len(cdrdata)
		for i := 0; i < cdrlen-1; i++ {
			_cdrdata := strings.Split(cdrdata[i], "_")
			cdrfold := _cdrdata[3][0:8]
			fold, _ := Exec("cd /root/eGW/CDR/cdrDat; ls -l | grep '^d' | grep " + cdrfold + "|awk '{print $9}'")
			//fmt.Println(len(fold))
			//fmt.Println(fold)
			if len(fold) == 1 {
				Run("mkdir -p /root/eGW/CDR/cdrDat/" + cdrfold)
			}
			Run("mv /root/eGW/CDR/" + cdrdata[i] + " /root/eGW/CDR/cdrDat/" + cdrfold)
		}
	}
}

func Manager() {
	/*cdr话单处理*/
	go Managecdr()
}
