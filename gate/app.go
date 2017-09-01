package gate

import (
	"fmt"
	"github.com/labstack/echo"
	"gogate/gate/facts"
	"gogate/gate/handlers"
	"gogate/gate/types"
	"gogate/gate/utils"
	"time"
)

var status types.State

var contents types.Content = make(types.Content)

func run_for_ever(options types.Opt) {
	for {
		looper(options)
		time.Sleep(5000 * time.Millisecond)
	}
}

func looper(options types.Opt) {
	status.Os = facts.GetOs()
	status.Hostname, _ = facts.GetHostname()
	status.Interfaces, _ = facts.GetIfs()
	if options.Is_master {
		key := utils.FindKey(status, options.Key, contents)
		contents[key] = status
	} else {
		utils.ForwardToMaster(options.Master_addr, status)
	}
}

func App(options types.Opt) {

	e := echo.New()
	// pass options to context
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ccc := &types.CustomContext{c, options, status, contents}
			return h(ccc)
		}
	})

	go run_for_ever(options)

	e.GET("/self", handlers.Self)
	e.GET("/", handlers.Index)
	e.POST("/", handlers.Create)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", options.Listen)))
}
