package facts

import (
	"io/ioutil"
	"runtime"
)

func GetMemory() string {
	switch os := runtime.GOOS; os {
	case "linux":
		return getMemoryLinux()
	case "default":
		return "Unknown"
	}
	return "Unkonwn"
}

func getMemoryLinux() string {
	if bytes, err := ioutil.ReadFile("/proc/meminfo"); err == nil {
		return string(bytes)
	}
	return ""
}
