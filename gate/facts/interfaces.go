//get all network interfaces
package facts

import (
	"bytes"
	"fmt"
	"github.com/woosley/gogate/gate/types"
	"net"
	"strings"
)

var blacklist []string = []string{"lo"}

func GetIfs() ([]types.Intf, error) {
	ifs, err := net.Interfaces()
	if err != nil {
		return []types.Intf{}, err
	}

	_ifs := make([]types.Intf, 0)
OUTER:
	for _, i := range ifs {
		name := i.Name

		//skip if not up
		if i.Flags&net.FlagUp == 0 {
			continue OUTER
		}
		//skip if no mac address
		if bytes.Compare(i.HardwareAddr, nil) == 0 {
			continue OUTER
		}

		mac := i.HardwareAddr.String()

		//skip blacklist
		for _, bl := range blacklist {
			if string(name) == bl {
				continue OUTER
			}
		}

		ips := make([]string, 0)
		adds, err := i.Addrs()
		if err == nil {
			for _, addr := range adds {
				ips = append(ips, fmt.Sprintf("%s", strings.Split(addr.String(), "/")[0]))
			}
		}
		_ifs = append(_ifs, types.Intf{name, mac, ips})
	}
	return _ifs, nil
}
