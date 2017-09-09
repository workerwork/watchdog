package task

import (
//"strconv"
)

//获取系统信息
func Getcpuall() []string {
	cpuall, _ := Exec("sar -P ALL 1 1 |grep 'Average' |awk '{print $8}'")
	len := len(cpuall) - 1
	cpuidle := cpuall[2:len]
	/*var cpuidle []float64
	for _, v := range _cpuidle {
		fmt.Println(v)
		vv, _ := strconv.ParseFloat(v, 64)
		cpuidle = append(cpuidle, vv)
	}*/
	//fmt.Println(cpuidle)
	return cpuidle
}

func Getcpu() (cpu, core string) {
	_cpu, _ := Exec("ps -aux|grep 'ltegwd 0 1'$ |awk '{print $3}'")
	_core, _ := Exec("lscpu |grep '^CPU(s):' |awk '{print $2}'")
	return _cpu[0], _core[0]
}

func Getmem() (mem, memsum string) {
	_mem, _ := Exec("ps -aux|grep 'ltegwd 0 1'$ |awk '{print $4}'")
	_memsum, _ := Exec("free -h |grep Mem |awk '{print $2}'")
	return _mem[0], _memsum[0]
}

func Getdisk() (disksum, diskfree, diskegw, diskcdr, disklog string) {
	_disksum, _ := Exec("df -l |grep 'centos-root' |awk '{print $2}'")
	_diskfree, _ := Exec("df -l |grep 'centos-root' |awk '{print $4}'")
	_diskegw, _ := Exec("du -s /root/eGW |awk '{print $1}'")
	_diskcdr, _ := Exec("du -s /root/eGW/CDR |awk '{print $1}'")
	_disklog, _ := Exec("du -s /root/eGW/Logs |awk '{print $1}'")
	return _disksum[0], _diskfree[0], _diskegw[0], _diskcdr[0], _disklog[0]
}

func Getlink() (enblink, epclink string) {
	_enblink, _ := Exec("netstat -anp |grep ltegwd |grep sctp |grep ESTABLISHED |wc -l")
	_epclink, _ := Exec("netstat -anp |grep ltegwd |grep sctp |grep CLOSE |wc -l")
	return _enblink[0], _epclink[0]
}

func Getenbue() (enbnum, uenum string) {
	_enbnum, _ := Exec("/root/eGW/lccmd show_sum |grep 'Home eNodeB' |awk -F ':' '{print $2}'")
	_uenum, _ := Exec("/root/eGW/lccmd show_sum |grep 'UE' |awk -F ':' '{print $2}'")
	return _enbnum[0], _uenum[0]
}

func Getservice() (nginx, redis, monitor string) {
	_nginx, _ := Exec("systemctl status nginx |grep 'Active' |awk  '{print $2$3 }'")
	_redis, _ := Exec("systemctl status redis |grep 'Active' |awk  '{print $2$3 }'")
	_monitor, _ := Exec("systemctl status monitor |grep 'Active' |awk  '{print $2$3 }'")
	return _nginx[0], _redis[0], _monitor[0]
}
