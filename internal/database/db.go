package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DataBase struct {
	Connection *gorm.DB
}

func (database *DataBase) InitDB() error {
	var err error
	err = godotenv.Load("dev.env")
	if err != nil {
		log.Fatalf("Error load .env: %v", err)
	}
	host := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=" + sslmode
	logrus.Infoln("dsn: ", dsn)
	database.Connection, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Error open DB connection: ", err)
		return err
	}
	return nil
}

func (database *DataBase) CloseDB() error {
	sqlDB, err := database.Connection.DB()
	if err != nil {
		log.Println("Error close DB connection: ", err)
		return err
	}
	sqlDB.Close()
	return nil
}
