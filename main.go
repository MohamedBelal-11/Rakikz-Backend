package main

import (
	"rakkiz-backend/database"
	"rakkiz-backend/models"
	"rakkiz-backend/router"

	"gorm.io/gorm"
)

func main() {
  db := database.ConnectDB()
  db.AutoMigrate(&models.User{})
  test(db)
  router.Runserver(db)
}



func test(db *gorm.DB) 