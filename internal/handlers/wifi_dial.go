package handlers

import (
	"github.com/mmadfox/rpi-wifi-web-manager/internal/linux"

	"github.com/gofiber/fiber/v2"
)

type wifiDialRequest struct {
	SSID      string `json:"ssid"`
	Password  string `json:"password"`
	SavePoint bool   `json:"savePoint"`
}

type okResponse struct {
	Ok bool `json:"ok"`
}

func WiFiDial(iface string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		req := new(wifiDialRequest)
		if err := ctx.BodyParser(req); err != nil {
			return err
		}
		err := linux.DialWiFi(req.SSID, req.Password, req.SavePoint, iface)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse{
				Error: err.Error(),
			})
		}
		return ctx.JSON(okResponse{Ok: true})
	}
}
