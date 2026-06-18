package database

import (
	"fmt"
	"rakkiz-backend/errors"
	"rakkiz-backend/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() *gorm.DB {
  dsn := "host=aws-1-eu-central-1.pooler.supabase.com user=postgres.vjwkdeugeelqvvnvzohy password=V3?pnJjq9-7&Tx/ dbname=postgres port=5432 sslmode=require"

  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Silent),
  })
  if err != nil {
    fmt.Println("failed to connect database")
    log.Erorr(errors.Errorize(err))
  } else {
	fmt.Println("Connected to database")
  }

  return db
}