package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func webserviceStart() {
	log.Println("Starting webservice ...")

	e := echo.New()

	r := e.Group("/film")

	r.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.Logger.Fatal(e.Start(":4000"))
}
