package handlers

import (
	"github.com/mmadfox/rpi-wifi-web-manager/internal/linux"

	"github.com/gofiber/fiber/v2"
)

type ifacesResponse struct {
	Ifaces []linux.IfaceInfo `json:"ifaces"`
}

func Ifaces() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ifaces, err := linux.GetInterfaces()
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse{
				Error: err.Error(),
			})
		}
		return ctx.JSON(ifacesResponse{
			Ifaces: ifaces,
		})
	}
}
