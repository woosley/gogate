package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/woosley/gogate/gate/types"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// find key value for types.State
func FindKey(body types.State, key string, contents *types.Content) string {
	blacklist := []string{"docker0", "lo"}
	if key == "ip" {
		var ip string
		var last_ipv4 string
		for _, v := range body.Interfaces {
			// if there is already a same key in content

			if ListHasString(blacklist, v.Name) {
				continue
			}

			for _, _ip := range v.Ips {
				_, exists := contents.Get(_ip)
				if exists {
					return _ip
				}
				if !strings.Contains(_ip, "::") {
					last_ipv4 = _ip
				}
				ip = _ip
			}
		}
		if last_ipv4 != "" {
			return last_ipv4
		}
		return ip
	}

	if key == "mac" {
		var mac string
		for _, v := range body.Interfaces {
			// if there is already a same key in content
			_, exists := contents.Get(v.Mac)
			if exists {
				return v.Mac
			}
			mac = v.Mac
		}
		return mac
	}

	//hostname as default
	return body.Hostname
}

func ForwardToMaster(master string, data types.State) (string, error) {
	content, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Post(master, "text/plain", bytes.NewReader(content))
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	if res.StatusCode != 201 {
		return "", errors.New(fmt.Sprintf("server returned %d, not 201/Created", res.StatusCode))
	}
	if body, err := ioutil.ReadAll(res.Body); err != nil {
		return "", err
	} else {
		return string(body), nil
	}
}

func IsDir(f string) (bool, error) {
	finfo, err := os.Stat(f)
	if err != nil {
		return false, err
	}

	switch {
	case finfo.IsDir():
		return true, nil
	default:
		return false, nil
	}
}

func IsFile(f string) (bool, error) {
	finfo, err := os.Stat(f)
	if err != nil {
		return false, err
	}

	switch finfo.Mode().IsRegular() {
	case true:
		return true, nil
	default:
		return false, nil
	}
}

func ListHasString(list []string, s string) bool {
	for _, v := range list {
		if s == v {
			return true
		}
	}
	return false
}
