package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"

	"./api/v1"
)

func Router() http.Handler {
	router := chi.NewRouter()
	router.Mount("/api/v1/", v1.BtcRouter())
	return router
}

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.LUTC | log.Lshortfile)

	if os.Getenv("MAIN_PORT") == "" {
		log.Fatalln("Main port is not set")
	}

	if os.Getenv("RPC_ADDR") == "" && os.Getenv("ENV") != "test" {
		log.Fatalln("RPC server address is not set")
	}

	log.Println("Server is starting on port", os.Getenv("MAIN_PORT"))

	err := http.ListenAndServe(":"+os.Getenv("MAIN_PORT"), Router())

	if err != nil {
		log.Fatalln("Failed to start server with error:", err)
	}
}
