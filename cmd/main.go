package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"
	_ "time/tzdata" // Часовые пояса доступны и в минимальном Docker-образе.

	"github.com/vlonebara/todo-list/pkg/config"
	"github.com/vlonebara/todo-list/pkg/db"
	"github.com/vlonebara/todo-list/pkg/handlers"
	"github.com/vlonebara/todo-list/pkg/server"
	"github.com/vlonebara/todo-list/pkg/service"
)

// main загружает конфигурацию и запускает сервер.
func main() {
	backup := flag.String("backup", "", "save a database backup to a new file")
	restore := flag.String("restore", "", "restore a backup into TODO_DBFILE (must not exist)")
	flag.Parse()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	if *backup != "" || *restore != "" {
		if err := databaseCommand(cfg, *backup, *restore); err != nil {
			log.Fatal(err)
		}
		return
	}

	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

// databaseCommand выполняет резервное копирование без запуска HTTP-сервера.
func databaseCommand(cfg config.Config, backup, restore string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if backup != "" && restore != "" {
		return fmt.Errorf("choose either backup or restore")
	}
	if restore != "" {
		return db.Restore(ctx, restore, cfg.DBFile)
	}
	database, err := db.OpenDatabase(cfg.DBFile)
	if err != nil {
		return err
	}
	defer database.Close()
	return db.Backup(ctx, database, backup)
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
	handler := handlers.New(tasks, authentication, cfg.SecureCookies, "./web")
	return server.Run(cfg.Host, cfg.Port, handler)
}
