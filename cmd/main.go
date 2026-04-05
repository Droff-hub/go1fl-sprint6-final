package main

import (
	"log"
	"os"

	"github.com/Yandex-Practicum/golf1-sprint6-final/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	srv := server.NewServer(logger)

	logger.Println("Запуск сервера на http://localhost:8080")

	err := srv.Server.ListenAndServe()
	if err != nil {
		logger.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}