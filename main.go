package main

import (
	"fmt"
	"jobs/controllers"
	"jobs/models"
	"log"
	"net/http"

	. "jobs/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	DBURI := fmt.Sprintf("root:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", Cfg.DBPassword, Cfg.DBHost, Cfg.DBName)
	log.Println("db_uri:", DBURI)
	conn, err := gorm.Open(mysql.Open(DBURI), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		log.Fatal("when connecting to db:", err)
	}
	models.Init(conn)
	models.Migrate()

	router := controllers.GetRouter()

	if err := http.ListenAndServe(Cfg.Host, router); err != nil {
		log.Printf("when serving: %v", err)
	}
}
