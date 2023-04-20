package core

import (
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/templari/shire-client/core/model"
)

type HttpServer struct {
	core     *Core
	listener net.Listener
	echo     *echo.Echo
}

func StartHttpServer(core *Core, listener net.Listener) {
	server := &HttpServer{
		core:     core,
		listener: listener,
	}
	server.Start()
}

func (s *HttpServer) Start() {
	e := echo.New()
	e.Listener = s.listener
	s.echo = e

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.POST("/message", func(c echo.Context) error {
		message := model.Message{}
		if err := c.Bind(&message); err != nil {
			return err
		}
		if err := s.core.ReceiveMessage(message); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, message)
	})

	e.Logger.Fatal(e.Start(""))
	defer e.Close()
}

func (s *HttpServer) Stop() {
	s.echo.Close()
}
