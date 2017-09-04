package facts

import (
	"bufio"
	"github.com/woosley/gogate/gate/types"
	"runtime"
	"strings"
)

func GetCpuFacts() (types.CpuInfo, error) {
	var cup types.CpuInfo
	switch os := runtime.GOOS; os {
	case "linux":
		return getCpuLinux()
	default:
		return cpu, nil
	}
}

func getCpuLinux() (types.CpuInfo, error) {
	var cup types.CpuInfo
	if fd, err := os.Open("/proc/cpuinfo"); err != nil {
		return cpu, err
	} else {
		cpuinfo := make(map[string]string)
		fscanner := bufio.NewScanner(fd)
		for fscanner.Scan() {
			kv := strings.SplitN(fscanner.Text(), ":", 2)
			cpuinfo[strings.Trim(kv[0], " ")] = strings.Trim(kv[1], " ")
		}
	}
}
