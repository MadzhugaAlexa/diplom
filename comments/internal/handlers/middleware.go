package handlers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
)

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		allParams := c.QueryParams()
		path := c.Path()
		method := c.Request().Method
		requiestID := c.QueryParam("request_id")

		msg := fmt.Sprintf("Request[%s] %s %s with params:", requiestID, method, path)
		for _, pname := range c.ParamNames() {
			msg = msg + fmt.Sprintf(" %s=%v", pname, c.Param(pname))
		}

		for name, value := range allParams {
			msg = msg + fmt.Sprintf(" %s=%v", name, value)
		}

		defer func(msg *string) {
			log.Print(*msg)
		}(&msg)

		err := next(c)
		if err != nil {
			msg = msg + " failed with " + err.Error()
			return err
		}

		status := c.Response().Status
		msg = msg + " status = " + strconv.Itoa(status)

		return nil
	}
}
