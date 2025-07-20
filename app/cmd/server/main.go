package main

import (
	"errors"
	"fmt"
	"goscreener/internal/app"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	/*	go func() {
		log.Println("Starting pprof on :6060")
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()*/

	var HttpListenPort = os.Getenv("HTTP_LISTEN_PORT")
	if HttpListenPort == "" {
		HttpListenPort = "8080"
	}

	var Hostname = os.Getenv("HOSTNAME")
	if Hostname == "" {
		Hostname = "localhost"
	}

	service, err := app.NewApp(Hostname, HttpListenPort)
	if err != nil {
		log.Fatal("{FATAL} ", err)
	}

	fmt.Println(fmt.Sprintf("Starting server on %s:%s", Hostname, HttpListenPort))
	err = service.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")
		return
	}
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
