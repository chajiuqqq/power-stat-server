package handler

import (
	"github.com/labstack/echo/v4"
	"mqtt-wx-forward/service"
)

type Handler struct {
	sv *service.Service
}

func NewHandler(sv *service.Service) *Handler {
	return &Handler{
		sv: sv,
	}
}

func (h *Handler) GetStat(c echo.Context) error {
	return nil

}
func (h *Handler) PostEnergyManagement(c echo.Context) error {
	return nil
}
