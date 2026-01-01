package database

import (
	"fmt"
	"log"
	"os"

	"tm-go-api/models"

	"gorm.io/driver/sqlserver" // <--- CHANGED FROM postgres
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func ConnectDb() {
	// dsn := fmt.Sprintf(
	// 	"host=%s user=%s password=%s database=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
	// 	os.Getenv("DB_HOST"),
	// 	os.Getenv("DB_USER"),
	// 	os.Getenv("DB_PASSWORD"),
	// 	os.Getenv("DB_NAME"),
	// 	os.Getenv("DB_PORT"),
	// )
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	server := os.Getenv("DB_HOST")
	port := 1433
	database := os.Getenv("DB_NAME") // <--- CRITICAL: Must match exactly

	// 2. Build the DSN
	// Note: 'encrypt=true' is required for Azure.
	// Note: 'database=' tells it NOT to look in master.
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&encrypt=true",
		username, password, server, port, database)

	// Note: We use sqlserver.Open() instead of postgres.Open()
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to Azure SQL. \n", err)
	}

	log.Println("connected to Azure SQL")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("running migrations")
	db.AutoMigrate(&models.User{})

	Database = DbInstance{
		Db: db,
	}
}
