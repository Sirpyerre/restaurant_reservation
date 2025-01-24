package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"restaurant_reservation/internal/database"
)

type HealthHandler struct {
	db database.Service
}

func NewHealthHandler(db database.Service) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

func (h *HealthHandler) HealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, h.db.Health())
}
