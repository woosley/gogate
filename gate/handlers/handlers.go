package handlers

import (
	"encoding/json"
	"fmt"
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
	ccc.Contents.Set(key, data)
	return ccc.JSON(http.StatusCreated,
		map[string]string{"message": "node created", "key": key, "code": fmt.Sprintf("%v", http.StatusCreated)})
}

func Health(c echo.Context) error {
	ccc := c.(*types.CustomContext)
	m := map[string]string{"state": "running"}
	return ccc.JSON(http.StatusOK, m)
}

func GetNode(c echo.Context) error {
	ccc := c.(*types.CustomContext)
	key := ccc.Param("key")
	v, ok := ccc.Contents.Get(key)
	if !ok {
		return ccc.JSON(http.StatusNotFound,
			map[string]string{"message": "node not found!", "code": fmt.Sprintf("%v", http.StatusNotFound)})
	} else {
		return ccc.JSON(http.StatusOK, v)
	}
}

func DeleteNode(c echo.Context) error {
	ccc := c.(*types.CustomContext)
	key := ccc.Param("key")
	_, ok := ccc.Contents.Get(key)
	if !ok {
		return ccc.JSON(http.StatusNotFound,
			map[string]string{"message": "node not found!", "code": fmt.Sprintf("%v", http.StatusNotFound)})
	} else {
		ccc.Contents.Delete(key)
		return ccc.JSON(http.StatusOK,
			map[string]string{"message": "node deleted!", "code": fmt.Sprintf("%v", http.StatusOK)})
	}
}
