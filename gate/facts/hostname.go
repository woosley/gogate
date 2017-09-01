package facts

import "os"

func GetHostname() (string, error) {
	return os.Hostname()
}
