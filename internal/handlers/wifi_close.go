package handlers

import (
	"github.com/mmadfox/rpi-wifi-web-manager/internal/linux"

	"github.com/gofiber/fiber/v2"
)

func WiFiClose(iface string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := linux.CloseWiFi(iface)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse{
				Error: err.Error(),
			})
		}
		return ctx.JSON(okResponse{Ok: true})
	}
}
