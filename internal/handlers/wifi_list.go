package handlers

import (
	"time"

	"github.com/mmadfox/rpi-wifi-web-manager/internal/linux"

	"github.com/gofiber/fiber/v2"
)

type wifiListCommandResponse struct {
	ActiveSSID string            `json:"active"`
	List       []linux.WiFiPoint `json:"list"`
	UpdatedAt  int64             `json:"updatedAt"`
}

func WiFiList() fiber.Handler {
	cmd := linux.NewWiFiCommand()
	return func(ctx *fiber.Ctx) error {
		points, err := linux.ScanWiFiPoints(cmd)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).
				JSON(errorResponse{
					Error: err.Error(),
				})
		}
		resp := wifiListCommandResponse{
			ActiveSSID: "None",
			List:       points,
		}
		for i := 0; i < len(points); i++ {
			if points[i].Active {
				resp.ActiveSSID = points[i].SSID
				break
			}
		}
		resp.UpdatedAt = time.Now().Unix()
		return ctx.JSON(resp)
	}
}
