package app

import (
	"context"
	"fmt"
	"github.com/urashav/event-hub/configs"
	http_handler "github.com/urashav/event-hub/internal/handler/http"
	"github.com/urashav/event-hub/internal/middleware"
	repository "github.com/urashav/event-hub/internal/repository/postgres"
	"github.com/urashav/event-hub/internal/service"
	"github.com/urashav/event-hub/pkg/auth"
	"github.com/urashav/event-hub/pkg/hasher"
	httputils "github.com/urashav/event-hub/pkg/httputilst"
	"github.com/urashav/event-hub/pkg/postgres"
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

	db, err := postgres.NewPostgresDB(dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("Error closing database connection:", err)
		}
	}()
	hasher := hasher.NewHasher()
	user_repo := repository.NewUserRepository(db)
	tokenManager := auth.NewTokenManager(cfg.JWT.SigningKey)
	authMiddleware := middleware.NewAuthMiddleware(tokenManager)

	commonMiddleware := middleware.Chain(
		middleware.Recover,
		middleware.Logger,
		middleware.CORS,
	)

	user_service := service.NewUserService(
		user_repo,
		hasher,
		tokenManager,
	)
	user_handler := http_handler.NewUserHandler(user_service)

	// Роуты
	mux := http.NewServeMux()
	mux.Handle("/api/v1/user/signup", http.HandlerFunc(user_handler.SignUp))
	mux.Handle("/api/v1/user/signin", http.HandlerFunc(user_handler.SignIn))

	protected := middleware.Chain(
		commonMiddleware,
		authMiddleware.RequireAuth)

	mux.Handle("/api/v1/protected", protected(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value("user_id").(int)
			response := map[string]interface{}{
				"message": "Protected resource",
				"user_id": userID,
			}
			httputils.SendSuccess(w, response, http.StatusOK)
		}),
	))

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
