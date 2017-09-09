package task

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

/****************************************************/
const (
	version      = "WCG_V100R001C00B030"
	ltegwd       = "a147f85d69fa49a2f8e2577b2b30c755"
	ltegw_sh     = "5cb284bfe1e292c8b540cf1ee3830b7b"
	lccmd        = "693619c038a41067ca1d294361a5ff0c"
	gtp_relay_ko = "b936bed6f2aa7ff11d65711eb68b6410"
	monitor      = "4e6cf25eabe7c4db5f34c45208ceb730"
	startAll_sh  = "f5bd024b52d6b201401021f5ec385930"
	vtysh        = "e46e2dc7d03ab88495e472ae9dfefe34"
	watchdog     = "a002ae0d5916fa66a54e6993f509b5cd"
	register     = "268843be3b27b6e4a6d9c06d7f7ae9c6"
	egw_manage   = "1eba0bdc16753f5f7ee8a7e4cbfc7b4d"
	egw_report   = "251ca567f93aa4c9c54982a91aab903e"
)

/***************************************************/

func Check() {
	c := make(chan []string)
	for {
		go func() {
			c <- readRedis("BRPOP", "eGWstat_check", "0")
		}()
		str := <-c
		switch str[1] {
		case "version":
			cmd := `ls -l /root/eGW |grep ltegwd$ |awk '{print $6 " " $7 " " $8}'
					ls -l /root/eGW |grep ltegw.sh$ |awk '{print $6 " " $7 " " $8}'
					ls -l /root/eGW |grep lccmd$ |awk '{print $6 " " $7 " " $8}'
					ls -l /root/eGW |grep gtp-relay.ko$ |awk '{print $6 " " $7 " " $8}'
					ls -l /root/eGW |grep monitor$ |awk '{print $6 " " $7 " " $8}'
					ls -l /root/eGW |grep startAll.sh$ |awk '{print $6 " " $7 " " $8}'
					ls -l /root/eGW |grep vtysh$ |awk '{print $6 " " $7 " " $8}'
					ls -l /root/eGW |grep watchdog$ |awk '{print $6 " " $7 " " $8}'
					ls -l /root/eGW/Licence |grep register$ |awk '{print $6 " " $7 " " $8}'
					ls -l /root/eGW/OMC |grep egw_manage$ |awk '{print $6 " " $7 " " $8}'
					ls -l /root/eGW/OMC |grep egw_report$ |awk '{print $6 " " $7 " " $8}'`

			result, _ := exec.Command("/bin/bash", "-c", cmd).Output()
			version_time := strings.Split(string(result), "\n")
			if existsRedis("eGWstat_version_time") == true {
				writeRedis("DEL", "eGWstat_version_time")
			}
			len_version_time := len(version_time)
			for i := 0; i < len_version_time-1; i++ {
				writeRedis("RPUSH", "eGWstat_version_time", version_time[i])
				//writeRedis("RPUSH","eGWstat_readflag", "OK")
			}
			cmd = `md5sum /root/eGW/ltegwd |awk '{print $1}';
				   mddd5sum /root/eGW/ltegw.sh |awk '{print $1}';
				   md5sum /root/eGW/lccmd |awk '{print $1}';
				   md5sum /root/eGW/gtp-relay.ko |awk '{print $1}';
				   md5sum /root/eGW/monitor |awk '{print $1}';
				   md5sum /root/eGW/startAll.sh |awk '{print $1}';
				   md5sum /root/eGW/vtysh |awk '{print $1}';
				   md5sum /root/eGW/watchdog |awk '{print $1}';
				   md5sum /root/eGW/Licence/register |awk '{print $1}';
				   md5sum /root/eGW/OMC/egw_manage |awk '{print $1}';
				   md5sum /root/eGW/OMC/egw_report |awk '{print $1}'`
			result, _ = exec.Command("/bin/bash", "-c", cmd).Output()
			md5sum := strings.Split(string(result), "\n")
			if existsRedis("eGWstat_version") == true {
				writeRedis("DEL", "eGWstat_version")
			}
			len_md5sum := len(md5sum)
			writeRedis("RPUSH", "eGWstat_version", version)
			for i := 0; i < len_md5sum-1; i++ {
				writeRedis("RPUSH", "eGWstat_version", md5sum[i])
			}

			if existsRedis("eGWstat_version_std") == true {
				writeRedis("DEL", "eGWstat_version_std")
			}
			writeRedis("RPUSH", "eGWstat_version_std", version, ltegwd, ltegw_sh, lccmd, gtp_relay_ko, monitor, startAll_sh, vtysh, watchdog, register, egw_manage, egw_report)
			if existsRedis("eGWstat_readflag") == true {
				writeRedis("DEL", "eGWstat_readflag")
			}
			writeRedis("RPUSH", "eGWstat_readflag", "OK")
		case "system":
			cpu, core := Getcpu()
			mem, memsum := Getmem()
			disksum, diskfree, diskegw, diskcdr, disklog := Getdisk()
			if existsRedis("eGWstat_system") == true {
				writeRedis("DEL", "eGWstat_system")
			}
			writeRedis("RPUSH", "eGWstat_system", cpu, core, mem, memsum, disksum, diskfree, diskegw, diskcdr, disklog)
			//判断阈值，是否上报
			_cpu, _ := strconv.ParseFloat(cpu, 32)
			_mem, _ := strconv.ParseFloat(mem, 32)
			_disksum, _ := strconv.Atoi(disksum)
			_diskegw, _ := strconv.Atoi(diskegw)
			_disksum_ := float32(_disksum)
			_diskegw_ := float32(_diskegw)

			cpuidle := Getcpuall()
			var _cpuall []string
			var _cpuflag bool = false
			//var vvv float64
			for _, v := range cpuidle {
				//fmt.Println(v)
				vv, _ := strconv.ParseFloat(v, 64)
				vvv := fmt.Sprintf("%.2f", 100-vv)
				//pow10_n := math.Pow10(2)
				//vvv := math.Trunc((100-vv+0.5/pow10_n)*pow10_n) / pow10_n
				_cpuall = append(_cpuall, vvv)
				//writeRedis("RPUSH", "eGWstat_cpucore", vvv)
				if 100.0-vv >= 80.0 {
					_cpuflag = true
				}
			}
			//fmt.Println(_cpuall)
			if existsRedis("eGWstat_cpucore") == true {
				writeRedis("DEL", "eGWstat_cpucore")
			}
			for _, v := range _cpuall {
				writeRedis("RPUSH", "eGWstat_cpucore", v)
			}

			if _cpu >= 80.0 || _mem >= 80.0 || _diskegw_/_disksum_ >= 0.8 || _cpuflag == true {
				if existsRedis("eGWstat_system_rp") == true {
					writeRedis("DEL", "eGWstat_system_rp")
				}
				writeRedis("RPUSH", "eGWstat_system_rp", cpu, core, mem, memsum, disksum, diskfree, diskegw)
				if existsRedis("eGWstat_cpucore_rp") == true {
					writeRedis("DEL", "eGWstat_cpucore_rp")
				}
				for _, v := range _cpuall {
					writeRedis("RPUSH", "eGWstat_cpucore_rp", v)
				}
			}
			if existsRedis("eGWstat_readflag") == true {
				writeRedis("DEL", "eGWstat_readflag")
			}
			writeRedis("RPUSH", "eGWstat_readflag", "OK")
		case "info":
			enblink, epclink := Getlink()
			enbnum, uenum := Getenbue()
			if existsRedis("eGWstat_info") == true {
				writeRedis("DEL", "eGWstat_info")
			}
			writeRedis("RPUSH", "eGWstat_info", enblink, epclink, enbnum, uenum)
			if existsRedis("eGWstat_readflag") == true {
				writeRedis("DEL", "eGWstat_readflag")
			}
			writeRedis("RPUSH", "eGWstat_readflag", "OK")
		case "service":
			nginx, redis, monitor := Getservice()
			if existsRedis("eGWstat_service") == true {
				writeRedis("DEL", "eGWstat_service")
			}
			writeRedis("RPUSH", "eGWstat_service", nginx, redis, monitor)
			if existsRedis("eGWstat_readflag") == true {
				writeRedis("DEL", "eGWstat_readflag")
			}
			writeRedis("RPUSH", "eGWstat_readflag", "OK")
			//time.Sleep(100 * time.Millisecond)
		}
	}
}

func Checker() {
	/*查询*/
	go Check()

}
