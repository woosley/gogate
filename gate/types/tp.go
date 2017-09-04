package types

import (
	"github.com/labstack/echo"
)

type Intf struct {
	Name string
	Mac  string
	Ips  []string
}

type CustomContext struct {
	echo.Context
	Opts     Opt
	Status   State
	Contents Content
}

type Opt struct {
	Listen      int
	Is_master   bool
	Help        bool
	Master_addr string
	Expire      int
	Key         string
}

type State struct {
	Os         string
	Hostname   string
	Interfaces []Intf
	Memory     string
	Swap       string
}

type Content map[string]State
