package types

import (
	"encoding/json"
	"github.com/labstack/echo"
	"sync"
)

type App struct {
	Name    string
	Health  HealthStatus
	Version string
}

type HealthStatus struct {
	Url    string
	State  string
	Code   int
	Reason string
}

type Intf struct {
	Name string
	Mac  string
	Ips  []string
}

type CustomContext struct {
	echo.Context
	Opts     Opt
	Status   State
	Contents *Content
}

type Opt struct {
	Listen      int
	Is_master   bool
	Help        bool
	Master_addr string
	Expire      int
	Key         string
	Version     bool
}

type State struct {
	Os         string
	Hostname   string
	Apps       []App
	Interfaces []Intf
	Memory     string
	Swap       string
	Cpu        CpuInfo
	Uptime     int64
	LastUpdate int64
}

type CpuInfo struct {
	Count int
	Cores int
}

type Content struct {
	sync.RWMutex
	m map[string]State
}

func NewContent() *Content {
	return &Content{
		m: make(map[string]State),
	}
}

func (c *Content) Get(key string) (State, bool) {
	c.RLock()
	defer c.RUnlock()
	value, ok := c.m[key]
	return value, ok
}

func (c *Content) Set(key string, value State) {
	c.Lock()
	defer c.Unlock()
	c.m[key] = value
}

func (c *Content) Delete(key string) {
	c.Lock()
	defer c.Unlock()
	delete(c.m, key)
}

func (c *Content) MarshalJSON() ([]byte, error) {
	c.RLock()
	defer c.RUnlock()
	return json.Marshal(c.m)
}

func (c *Content) UnMarshalJSON(data []byte) error {
	c.Lock()
	defer c.Unlock()
	err := json.Unmarshal(data, &c.m)
	return err
}
