package main

import (
	"fmt"
	"rakkiz-backend/database"
	"rakkiz-backend/log"
	"rakkiz-backend/models"
	"rakkiz-backend/slices"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
  db := database.ConnectDB()
  db.AutoMigrate(&models.User{})
  test(db)
}

func runserver() {
  r := gin.Default()
  r.GET("/", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "Hello World",
    })
  })
  r.Run(":8080")
}

func test(db *gorm.DB) {
  userService := models.UserService{Db: db}
  user, err := userService.Create(&models.User{
    Username: "Mohamed-Belal",
    Email: "mo7amedbll@gmail.com",
    Name: "Mohamed Belal",
    Password: "Secret123!",
    IsVerified: true,
  })

  if err != nil {
    log.Erorr(err)
  } else {
    fmt.Println("User created: ", user.Username)
  }

  users := userService.Objects()
  fmt.Println("Users:\n", strings.Join(
    slices.Map(users.All, func(u models.User) string {
      return fmt.Sprintf("- %s (%s) %s %s", u.Username, u.Email, u.Name, u.CreatedAt.Format("2006-01-02 15:04:05"))
    }),
    "\n"))
}