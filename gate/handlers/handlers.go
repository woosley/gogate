package handlers

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/woosley/gogate/gate/types"
	"github.com/woosley/gogate/gate/utils"
	"io/ioutil"
	"net/http"
)

func Self(c echo.Context) error {
	ccc := c.(*types.CustomContext)

	return ccc.JSON(http.StatusOK, ccc.Status)
}

func Index(c echo.Context) error {
	ccc := c.(*types.CustomContext)
	// if master, print result
	if ccc.Opts.Is_master {
		return ccc.JSON(http.StatusOK, ccc.Contents)
		// if not master, redirect to master address
	} else {
		return ccc.Redirect(http.StatusMovedPermanently, ccc.Opts.Master_addr)
	}
}

// Post to / to create a new recode
func Create(c echo.Context) error {
	ccc := c.(*types.CustomContext)
	if !ccc.Opts.Is_master {
		m := make(map[string]string)
		m["reason"] = "Can not post to nodes"
		return ccc.JSON(http.StatusBadRequest, m)
	}

	var data types.State
	body, _ := ioutil.ReadAll(ccc.Request().Body)

	err := json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	key := utils.FindKey(data, ccc.Opts.Key, ccc.Contents)
	ccc.Contents[key] = data
	return ccc.String(http.StatusCreated, key)
}

func Health(c echo.Context) error {
	ccc := c.(*types.CustomContext)
	m := map[string]string{"state": "running"}
	return ccc.JSON(http.StatusOK, m)
}
