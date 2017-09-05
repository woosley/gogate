package facts

import (
	"bufio"
	"github.com/woosley/gogate/gate/types"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func GetCpu() (types.CpuInfo, error) {
	var cpu types.CpuInfo
	switch os := runtime.GOOS; os {
	case "linux":
		return getCpuLinux()
	default:
		return cpu, nil
	}
}

func getCpuLinux() (types.CpuInfo, error) {
	var cpu types.CpuInfo
	if fd, err := os.Open("/proc/cpuinfo"); err != nil {
		return cpu, err
	} else {
		cpuinfo := make(map[string]string)
		fscanner := bufio.NewScanner(fd)
		for fscanner.Scan() {
			line := fscanner.Text()
			if strings.Contains(line, ":") {
				kv := strings.SplitN(fscanner.Text(), ":", 2)
				cpuinfo[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			}
		}
		cpu.Count, _ = strconv.Atoi(cpuinfo["physical id"])
		cpu.Count += 1
		cpu.Cores, _ = strconv.Atoi(cpuinfo["cpu cores"])
		return cpu, nil
	}
}
