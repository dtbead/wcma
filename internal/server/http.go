package server

import (
	"context"

	"github.com/dtbead/wc-maps-archive/internal/service"
	"github.com/labstack/echo/v4"
)

type Message struct {
	Error string `json:"error"`
}

type ServerController struct {
	e          *echo.Echo
	service    *service.Service
	videoGroup *echo.Group
}

func NewServer(s *service.Service) ServerController {
	ctrl := ServerController{
		echo.New(),
		s, nil}

	ctrl.videoGroup = ctrl.e.Group("/video")

	ctrl.initEcho()
	return ctrl
}

func (s ServerController) Start(address string) error {
	return s.e.Start(address)
}

func (s ServerController) Stop() error {
	return s.e.Shutdown(context.Background())
}

func (s ServerController) initEcho() {
	s.getVideoInfo()
}

func (s ServerController) getVideoInfo() {
}
