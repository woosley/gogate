// +build windows
package facts

import (
	"fmt"
	"github.com/StackExchange/wmi"
)

type WinOS struct {
	Version string
	Caption string
}

func GetOs() string {
	var dst WinOS
	q := wmi.CreateQuery(&dst, "")
	err := wmi.Query(q, &dst)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s %s", dst[0].Caption, dst[0].Version)
}
