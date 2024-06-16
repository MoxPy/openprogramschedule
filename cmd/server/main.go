package main

import (
	"context"
	"errors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"openprogramschedule/api/handlers"
	"openprogramschedule/api/routes"
	"openprogramschedule/internal/db"
	"openprogramschedule/internal/middlewares"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	database := db.ConnectDB()

	programEnv := &handlers.ProgramHandler{
		Db: database,
	}
	scheduleEnv := &handlers.ScheduleHandler{
		Db: database,
	}
	defer func() {
		err := db.CloseDB()
		if err != nil {
			log.Fatal(err)
		}
	}()

	mux := http.NewServeMux()
	routes.ProgramRouter(mux, programEnv)
	routes.ScheduleRouter(mux, scheduleEnv)

	wrappedMux := middlewares.AuthMiddleware(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: wrappedMux,
	}
	// Start
	go func() {
		log.Println("Server listening on port 8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-interruptChan
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
