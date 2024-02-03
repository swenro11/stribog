package app

import (
	"log"

	"github.com/swenro11/stribog/internal/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	databaseURL := "postgres://stribog:stribog@127.0.0.1:5432/stribog"
	databaseURL += "?sslmode=disable"
	//dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("gorm.Open error: %s", err)
	}
	db.AutoMigrate(&entity.Article{})
	log.Print("AutoMigrate: entity.Article")
	db.AutoMigrate(&entity.Image{})
	log.Print("AutoMigrate: entity.Image")
	db.AutoMigrate(&entity.Keyword{})
	log.Print("AutoMigrate: entity.Keyword")
}
