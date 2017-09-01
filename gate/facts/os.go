package facts

import (
	"io/ioutil"
	"os"
)

const redhat string = "/etc/redhat-release"

func GetOs() string {
	if _, err := os.Stat(redhat); !os.IsExist(err) {
		if bytes, err := ioutil.ReadFile("/etc/redhat-release"); err == nil {
			return string(bytes)
		}
	}
	return "Unknown"
}
