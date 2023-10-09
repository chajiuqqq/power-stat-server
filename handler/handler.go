package handler

import "mqtt-wx-forward/service"

type Handler struct {
	sv *service.Service
}

func NewHandler(sv *service.Service) *Handler {
	return &Handler{
		sv: sv,
	}
}
