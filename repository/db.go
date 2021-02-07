package repository

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Migrate() {
	db, err := gorm.Open(sqlite.Open("company-funding.db"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&CompanyFunding{})
}

func Connect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("company-funding.db"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	return db
}
