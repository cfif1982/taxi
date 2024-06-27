package main

import (
	"log"

	"github.com/cfif1982/taxi/internal"
	"github.com/cfif1982/taxi/pkg/logger"
)

func main() {

	// инициализируем логгер
	logger, err := logger.GetLogger()

	// Если логгер не инициализировался, то выводим сообщение с помощью обычного log
	if err != nil {
		log.Fatal("logger zap initialization: FAILURE")
	}

	// выводим сообщенеи об успешной инициализации логгера
	logger.Info("logger zap initialization: SUCCESS")

	// создаем сервер
	srv := internal.NewServer(logger)

	// запускаем сервер
	if err := srv.Run("localhost:8080"); err != nil {
		logger.Fatal("error occured while running http server: " + err.Error())
	}

}
