package facts

import (
	"bufio"
	"os"
	"runtime"
	"strings"
)

func GetMemory() (string, string, error) {
	switch os := runtime.GOOS; os {
	case "linux":
		return getMemoryLinux()
	default:
		return "", "", nil
	}
}

//getMemoryLinux returns both MemTotal and SwapTotal
func getMemoryLinux() (string, string, error) {
	if fd, err := os.Open("/proc/meminfo"); err != nil {
		return "", "", err
	} else {
		mem := make(map[string]string)
		fscanner := bufio.NewScanner(fd)
		for fscanner.Scan() {
			kv := strings.SplitN(fscanner.Text(), ":", 2)
			mem[kv[0]] = strings.Trim(kv[1], " ")
		}
		return mem["MemTotal"], mem["SwapTotal"], nil
	}
}
