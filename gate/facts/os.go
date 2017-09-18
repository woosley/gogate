package facts

import (
	"bufio"
	"github.com/woosley/gogate/gate/utils"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

const unknown string = "Unknown"

var osfiles []string = []string{
	"/etc/redhat-release",
	"/etc/os-release",
}

func GetOs() string {
	switch osf := runtime.GOOS; osf {
	case "linux":
		return getLinuxDistro()
	case "darwin":
		return getMacDistro()
	default:
		return unknown
	}
}

func getOsFile(osfs []string) string {
	for _, f := range osfs {
		if exists, err := utils.IsFile(f); err == nil && exists {
			return f
		}
	}
	return ""
}

func getLinuxDistro() string {

	switch f := getOsFile(osfiles); f {
	case "/etc/redhat-release":
		if bytes, err := ioutil.ReadFile("/etc/redhat-release"); err == nil {
			return strings.TrimSpace(string(bytes))
		}
	case "/etc/os-release":
		return getFromOsRelease(f)

	default:
		return unknown
	}
	return unknown
}

//getFromOsRelease return distribution name from /etc/os-release
func getFromOsRelease(f string) string {
	fd, err := os.Open(f)
	if err != nil {
		return unknown
	}

	fscanner := bufio.NewScanner(fd)

	for fscanner.Scan() {
		line := fscanner.Text()
		if strings.Contains(line, "=") {
			kv := strings.SplitN(line, "=", 2)
			k, v := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
			if k == "PRETTY_NAME" {
				return strings.Trim(v, "\"")
			}
		}
	}
	return unknown
}

func getMacDistro() string {
	return "MacOS"
}
