package models

import (
	"gorm.io/gorm"
)

var db *gorm.DB

// GetDB returns db
func GetDB() *gorm.DB {
	return db
}

// Init db using by models with conn
func Init(conn *gorm.DB) {
	db = conn
}

var activeModels = []interface{}{&Job{}, &Task{}, &Work{}}

// Migrate (with debug to see all queries in logs)
func Migrate() {
	if err := GetDB().Debug().AutoMigrate(activeModels...); err != nil {
		panic("when tryin to migrate: " + err.Error())
	}
}

//Truncate (handy when testing)
func Truncate() {
	for _, m := range activeModels {
		GetDB().Unscoped().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(m)
	}
}
