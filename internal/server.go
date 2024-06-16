package internal

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cfif1982/taxi/internal/application/middlewares"
	"github.com/cfif1982/taxi/pkg/logger"
	"github.com/go-chi/chi/v5"

	driversHandler "github.com/cfif1982/taxi/internal/application/drivers/handlers"
	routesHandler "github.com/cfif1982/taxi/internal/application/routes/handlers"

	driversInfra "github.com/cfif1982/taxi/internal/infrastructure/drivers"
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
		s.logger.Fatal("can't initialize postgres for routes DB: " + err.Error())
	} else {
		s.logger.Info("postgres for routes DB initialized")
	}

	// Создаем репозиторий для работы с водителями
	driverRepo, err := driversInfra.NewPostgresRepository(ctx, databaseDSN, s.logger)

	if err != nil {
		s.logger.Fatal("can't initialize postgres for drivers DB: " + err.Error())
	} else {
		s.logger.Info("postgres for drivers DB initialized")
	}

	// создаем хэндлер маршрута и передаем ему нужную БД
	routeHandler := routesHandler.NewHandler(routeRepo, s.logger)

	// создаем хэндлер водителя
	driverHandler := driversHandler.NewHandler(driverRepo, s.logger)

	// создаем роутер
	routerChi := chi.NewRouter()

	// назначаем middleware
	routerChi.Use(middlewares.GzipCompressMiddleware)
	routerChi.Use(middlewares.GzipDecompressMiddleware)

	// назначаем ручки для водителя
	s.SetDriverHandlers(routerChi, driverHandler)

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
	router.Get(`/api/admin/route`, middlewares.AdminAuthMiddleware(handler.GetRoute()))
	router.Put(`/api/admin/route`, middlewares.AdminAuthMiddleware(handler.EditRoute()))
	router.Delete(`/api/admin/route`, middlewares.AdminAuthMiddleware(handler.DeleteRoute()))
}

// назначаем ручки для водителя
func (s *Server) SetDriverHandlers(router *chi.Mux, handler *driversHandler.Handler) {

	// router.Get(`/api/driver/start`, handler.Start())
	router.Post(`/api/driver/registration`, handler.Registration())
	router.Post(`/api/driver/login`, handler.DriverLogin())
	router.Put(`/api/driver/balance`, handler.IncreaseBalance())
}
