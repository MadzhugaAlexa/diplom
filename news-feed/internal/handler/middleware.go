package handler

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
		requestID := c.QueryParam("request_id")
		ips := c.Request().Header["X-Forwarded-For"]
		var ip string
		if len(ips) > 0 {
			ip = ips[0]
		}

		msg := fmt.Sprintf("[IP: %s] Request[%s] %s %s with params:", ip, requestID, method, path)
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
		msg = msg + " код ответа = " + strconv.Itoa(status)

		return nil
	}
}
