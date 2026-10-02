package main

import (
	"ad/make_server"
	"ad/postgre"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	server := make_server.Create()

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error : %s\n", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM)

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error during graceful shutdown : %s\n", err)
	}
	if err := postgre.Close(); err != nil {
		log.Fatalf("Error during graceful shutdown : %s\n", err)
	}
	log.Printf("The server is shutdown\n")
}
