package app

import (
	"context"
	"fmt"
	"github.com/urashav/event-hub/configs"
	"github.com/urashav/event-hub/database"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func App(cfg *configs.Config) {
	// Создаем канал для перехвата сигналов завершения, затем ждем сигнала
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// Устанавливаем таймаут для завершения работы сервера
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Подключение к базе данных
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name)

	db, err := database.NewPostgresDB(dsn)
	if err != nil {
		log.Println("Error connecting to database:", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("Error closing database connection:", err)
		}
	}()

	// Роуты
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!!!")
	})

	// Скервер запускаем в горутине, чтобы освободить основной поток
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: mux,
	}
	go func() {
		log.Println("Listening and serving HTTP on ", srv.Addr)
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println(err)
			quit <- syscall.SIGTERM // Отправляем сигнал завершения, если сервер не смог запуститься
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
