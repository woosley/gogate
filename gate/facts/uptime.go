package facts

import (
	"runtime"
	"strings"
	"strconv"
	"io/ioutil"
)

const zero int64 = 0

func GetUptime() int64 {
	switch osf := runtime.GOOS; osf {
	case "linux":
		return getUptimeLinux()
	default:
		return zero
	}
}

func getUptimeLinux() int64 {
	fname := "/proc/uptime"	
	if bytes, err := ioutil.ReadFile(fname); err == nil {
		up := strings.SplitN(strings.TrimSpace(string(bytes)), " ", 2)[0]
		upt, err := strconv.ParseFloat(up, 64)
		if err != nil {
			return zero
		}else{
			return int64(upt)
		}
	}

	return zero
}
