package main

import (
	"fmt"
	"jobs/controllers"
	"jobs/models"
	"log"
	"net/http"
	"time"

	"jobs/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// инициализируем базу
	cfg := config.GetConfig()
	dbURI := fmt.Sprintf("root:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.DBPassword, cfg.DBHost, cfg.DBName)
	log.Println("db_uri:", dbURI)
	// закомментированная опция нужна при разработке когда схема меняется
	conn, err := gorm.Open(mysql.Open(dbURI) /*&gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}*/)
	if err != nil {
		log.Fatal("when connecting to db:", err)
	}
	models.Init(conn)

	// при запуске актуализируем схему базы
	models.Migrate()

	// создаем сервер
	server := &http.Server{
		Handler: controllers.GetRouter(),
		Addr:    cfg.Host,
		// ставим таймауты чтобы запрос заканчивался если сервер завис
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("serving on", cfg.Host)
	// запускаем приложение
	if server.ListenAndServe(); err != nil {
		log.Printf("when serving: %v", err)
	}
}
