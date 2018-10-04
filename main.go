package main

import (
	"context"
	"fmt"
	"github.com/jdelgad/simpleGoHttpServer/internal/httpServer"
	"github.com/jdelgad/simpleGoHttpServer/internal/httpServer/handlers"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func main() {
	shutdownCh := make(chan bool)

	portStr := os.Getenv("HTTP_LISTEN_PORT")
	port, err := strconv.Atoi(portStr)

	listenAddress := fmt.Sprintf(":%v", port)

	s := httpServer.NewServer(listenAddress)

	log.Println("Initializing routes...")
	h := handlers.NewHandlers(shutdownCh)
	s.Handle(`/hash/?$`, h.Hash)
	s.Handle(`/hash/(\d+)/?$`, h.HashID)
	s.Handle(`/stats/?$`, h.Stats)

	s.Handle(`/shutdown/?$`, h.Shutdown)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		log.Println("Starting http server")

		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Println("ListenAndServe error:", err)
			close(shutdownCh)
		}
		log.Println("ListenAndServe completed")
	}()

	log.Println("Waiting for POST to /shutdown")
	<-shutdownCh

	log.Println("Received POST to /shutdown. Shutting down...")
	ctx := context.Background()
	err = s.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Error occurred while shutting down http server: %v", err)
	} else {
		log.Println("HTTP server successfully shutdown")
	}
	wg.Wait()
}
