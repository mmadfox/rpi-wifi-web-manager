package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mmadfox/rpi-wifi-web-manager/internal/handlers"
	"github.com/mmadfox/rpi-wifi-web-manager/ui"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

func main() {
	var (
		wifiIface = flag.String("wifi.iface", "wlan0", "Name of the wireless interface to use for Wi-Fi")
		httpHost  = flag.String("http.host", "0.0.0.0", "Host address for the server to listen on")
		httpPort  = flag.Int("http.port", 3003, "Port number for the server to listen on")
	)
	flag.Parse()

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		// TODO: handle errors, etc...
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/api/wifi-list", handlers.WiFiList())
	app.Post("/api/wifi-conn", handlers.WiFiDial(*wifiIface))
	app.Post("/api/wifi-close", handlers.WiFiClose())

	app.Get("/api/ifaces", handlers.Ifaces())
	app.Post("/api/ifaces/switch", handlers.SwitchIface())

	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(ui.StaticFiles),
		Index:      "index.html",
		PathPrefix: "src",
	}))

	go func() {
		addr := fmt.Sprintf("%s:%d", *httpHost, *httpPort)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Error starting server: %v\n", err)
		}
	}()

	fmt.Println("WebNM runnig...")
	terminate(app)
}

func terminate(app *fiber.App) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
	<-sigint
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	fmt.Println("WebNM closed.")
}
