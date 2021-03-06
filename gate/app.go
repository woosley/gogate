package gate

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/woosley/gogate/gate/facts"
	"github.com/woosley/gogate/gate/handlers"
	"github.com/woosley/gogate/gate/types"
	"github.com/woosley/gogate/gate/utils"
	"time"
)

var status types.State

var contents *types.Content = types.NewContent()

func runForever(options types.Opt) {
	for {
		looper(options)
		time.Sleep(1000 * time.Millisecond)
	}
}

func looper(options types.Opt) {
	status.Os = facts.GetOs()
	status.Hostname, _ = facts.GetHostname()
	status.Interfaces, _ = facts.GetIfs()
	status.Memory, status.Swap, _ = facts.GetMemory()
	status.Cpu, _ = facts.GetCpu()
	status.Uptime = facts.GetUptime()
	status.Disks = facts.GetDiskInfo()
	status.LastUpdate = time.Now().Unix()
	if options.Is_master {
		key := utils.FindKey(status, options.Key, contents)
		contents.Set(key, status)
	} else {
		utils.ForwardToMaster(options.Master_addr, status)
	}
}

func App(options types.Opt) {

	e := echo.New()
	e.Use(middleware.Logger())
	// pass options to context
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ccc := &types.CustomContext{c, options, status, contents}
			return h(ccc)
		}
	})

	if options.Debug {
		e.Logger.SetLevel(log.DEBUG)
	} else {
		e.Logger.SetLevel(log.INFO)
	}

	go runForever(options)
	e.HideBanner = true
	e.GET("/", handlers.Index)
	e.GET("/self", handlers.Self)
	e.GET("/health", handlers.Health)
	e.GET("/:key", handlers.GetNode)
	e.DELETE("/:key", handlers.DeleteNode)
	e.POST("/", handlers.Create)
	e.Logger.Info(fmt.Sprintf("Starting gogate on %v", options.Listen))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", options.Listen)))
}
