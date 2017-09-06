package facts

import (
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

const redhat string = "/etc/redhat-release"

func GetOs() string {
	switch osf := runtime.GOOS; osf {
	case "linux":
		return getLinuxDistro()
	case "darwin":
		return getMacDistro()
	default:
		return "Unknown"
	}
}

func getLinuxDistro() string {

	if _, err := os.Stat(redhat); !os.IsExist(err) {
		if bytes, err := ioutil.ReadFile("/etc/redhat-release"); err == nil {
			return strings.TrimSpace(string(bytes))
		}
	}
	return "Unknown"
}

func getMacDistro() string {
	return "MacOS"
}
