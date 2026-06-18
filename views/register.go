package views

import (
	"net/http"
	"rakkiz-backend/errors"
	"rakkiz-backend/models"

	"gorm.io/gorm"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	IsMuslim bool   `json:"is_muslim"`
}

func Register(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest

		if !isValidData(r.Body, req) {
			errorResponse(
				w,
				&errors.AppError{
					Code: 11017,
					Message: "Not valid data",
				},
				http.StatusBadRequest,
		  )
			return
		}

		userService := models.UserService{Db: db}
		_, err := userService.Create(&models.User{
			Username: req.Username,
			Email: req.Email,
			Password: req.Password,
			IsMuslim: &req.IsMuslim,
			IsVerified: true,
		})

		if err != nil {
			errorResponse(
				w,
				err,
				http.StatusBadRequest,
			)
			return
		}

		response(
			w,
			map[string]any{
				"success": true,
			},
			http.StatusCreated,
		)
	}
}
