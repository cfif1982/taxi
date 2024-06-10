package internal

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cfif1982/taxi/internal/application/middlewares"
	"github.com/cfif1982/taxi/pkg/logger"
	"github.com/go-chi/chi/v5"

	routesHandler "github.com/cfif1982/taxi/internal/application/routes/handlers"
	routesInfra "github.com/cfif1982/taxi/internal/infrastructure/routes"
)

// структура сервера
type Server struct {
	logger *logger.Logger
}

// Конструктор Server
func NewServer(logger *logger.Logger) Server {
	return Server{
		logger: logger,
	}
}

// запуск сервера
func (s *Server) Run(serverAddr string) error {

	// DSN для СУБД
	databaseDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", `localhost`, `postgres`, `123`, `taxi`)

	// создаю контекст для подключения БД
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Создаем репозиторий для работы с маршрутами
	routeRepo, err := routesInfra.NewPostgresRepository(ctx, databaseDSN, s.logger)

	if err != nil {
		s.logger.Fatal("can't initialize postgres route DB: " + err.Error())
	} else {
		s.logger.Info("postgres route DB initialized")
	}

	// // Создаем репозиторий для работы с водителями
	// driverRepo, err := driversInfra.NewPostgresRepository(ctx, databaseDSN, s.logger)

	// if err != nil {
	// 	s.logger.Fatal("can't initialize postgres driver DB: " + err.Error())
	// } else {
	// 	s.logger.Info("postgres driver DB initialized")
	// }

	// создаем хндлер маршрута и передаем ему нужную БД
	routeHandler := routesHandler.NewHandler(routeRepo, s.logger)
	//********************************************************

	// создаем роутер
	routerChi := chi.NewRouter()

	// назначаем middleware
	routerChi.Use(middlewares.GzipCompressMiddleware)
	routerChi.Use(middlewares.GzipDecompressMiddleware)

	// назначаем ручки для admin
	s.SetAdminHandlers(routerChi, routeHandler)

	s.logger.Info("Starting server", "addr", serverAddr)

	// запуск сервера на нужно адресе и с нужным роутером
	return http.ListenAndServe(serverAddr, routerChi)
}

// назначаем ручки для admin
func (s *Server) SetAdminHandlers(router *chi.Mux, handler *routesHandler.Handler) {

	router.Post(`/api/admin/login`, handler.AdminLogin())
	router.Get(`/api/admin/routes`, middlewares.AdminAuthMiddleware(handler.GetAllRoutes()))
	router.Post(`/api/admin/route`, middlewares.AdminAuthMiddleware(handler.AddRoute()))
}
