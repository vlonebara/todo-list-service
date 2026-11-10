package main

import (
	"log"

	"github.com/vlonebara/todo-list/pkg/config"
	"github.com/vlonebara/todo-list/pkg/db"
	"github.com/vlonebara/todo-list/pkg/handlers"
	"github.com/vlonebara/todo-list/pkg/server"
	"github.com/vlonebara/todo-list/pkg/service"
)

// main загружает конфигурацию и запускает сервер.
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

// run собирает зависимости приложения и закрывает БД после остановки сервера.
func run(cfg config.Config) error {
	authentication, err := service.NewAuthService(cfg.Password, cfg.JWTSecret)
	if err != nil {
		return err
	}
	database, err := db.OpenDatabase(cfg.DBFile)
	if err != nil {
		return err
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("close database: %v", err)
		}
	}()

	store := db.NewTaskStore(database)
	tasks := service.NewTaskService(store)
	handler := handlers.New(tasks, authentication, cfg.SecureCookies)
	return server.Run(cfg.Host, cfg.Port, handler)
}
