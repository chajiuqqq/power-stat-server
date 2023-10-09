package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"mqtt-wx-forward/handler"
	"mqtt-wx-forward/service"
	"mqtt-wx-forward/types"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	conf := types.NewConfig()
	opt := &service.ServiceOption{}
	logger := log.Default()
	sv := service.New(conf, logger, opt)
	h := handler.NewHandler(sv)
	// Routes
	pub := e.Group("/api")
	pub.GET("/api/wx/echo", h.GetWxEcho)
	pub.POST("/api/wx/echo", h.PostWxEcho)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
