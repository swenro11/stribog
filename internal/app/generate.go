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

	db.AutoMigrate(&entity.Task{})
	log.Print("AutoMigrate: entity.Task")
}
