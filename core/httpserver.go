package core

import (
	"net"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/templari/shire-client/model"
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
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, message)
	})

	e.POST("/groups/:id/prepare", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return err
		}
		err = s.core.PrepareGroup(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, "ok")
	})

	e.POST("/groups/:id/start", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return err
		}
		err = s.core.StartGroup(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, "ok")
	})

	e.Logger.Fatal(e.Start(""))
	defer e.Close()
}

func (s *HttpServer) Stop() {
	s.echo.Close()
}
