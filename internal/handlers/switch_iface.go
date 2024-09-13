package handlers

import (
	"github.com/mmadfox/rpi-wifi-web-manager/internal/linux"

	"github.com/gofiber/fiber/v2"
)

type switchIfaceRequest struct {
	Ifname string `json:"ifname"`
	Metric int    `json:"metric"`
}

func SwitchIface() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		req := new(switchIfaceRequest)
		if err := ctx.BodyParser(req); err != nil {
			return err
		}
		err := linux.SwitchInterface(req.Ifname, req.Metric)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse{
				Error: err.Error(),
			})
		}
		return ctx.JSON(okResponse{Ok: true})
	}
}
