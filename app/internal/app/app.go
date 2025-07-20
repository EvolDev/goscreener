package app

import (
	"fmt"
	"goscreener/internal/handlers/fakenav"
	"goscreener/internal/handlers/index"
	"goscreener/internal/handlers/screenshot"
	"net/http"
	"os"
	"time"
)

const addressPort = ":8080"

type App struct {
	mux    *http.ServeMux
	server *http.Server
}

func NewApp(host string, port string) (*App, error) {

	screenshotsDir := GetScreenshotsDir()
	fmt.Println("The screens are located in", screenshotsDir)

	screensDir := screenshotsDir
	if err := os.MkdirAll(screensDir, 0755); err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	mux.Handle("/", index.NewIndexHandler())

	// Get site screenshot by url
	mux.Handle("/screenshot", screenshot.NewScreenshotHandler())

	mux.Handle("/screenshots-many", screenshot.NewScreenshotsManyHandler(host, port, screensDir))

	// Get html for fake nav
	mux.Handle("/fake-nav", fakenav.NewFakeNavHandler())

	// Static files dir
	mux.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./htdocs/img"))))

	// Directory for screenshots saving via multiple requests
	mux.Handle("/get/screenshot/", http.StripPrefix("/get/screenshot/", http.FileServer(http.Dir(screensDir))))

	server := &http.Server{
		Addr:        addressPort,
		Handler:     mux,
		ReadTimeout: 5 * time.Second,
		IdleTimeout: 5 * time.Second,
	}

	return &App{
		mux:    mux,
		server: server,
	}, nil
}

func (a *App) ListenAndServe() error {
	return a.server.ListenAndServe()
}

func (a *App) Close() error {
	return a.server.Close()
}

func GetScreenshotsDir() string {
	dir := os.Getenv("SCREENSHOTS_DIR")
	if dir == "" {
		if _, err := os.Stat("/.env"); err == nil {
			// Running in Docker container
			dir = "/app/screens"
		} else {
			// Running in Local
			dir = "./screens"
		}
	}

	return dir
}
