// main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aasoft24/golara/bootstrap"
	"github.com/aasoft24/golara/wpkg/configs"
)

func main() {
	startServer()
}

func startServer() {
	router := bootstrap.Init() // router now returned from Init()

	host := configs.GConfig.Server.Host
	port := configs.GConfig.Server.Port

	serverAddr := fmt.Sprintf("%s:%d", host, port)
	url := fmt.Sprintf("http://%s", serverAddr)
	log.Printf("🚀 Server running at %s", url)

	// Use the returned router
	log.Fatal(http.ListenAndServe(serverAddr, router))
}
