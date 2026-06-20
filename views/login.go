package views

import (
	"net/http"
	"rakkiz-backend/errors"
	"rakkiz-backend/models"
	"strings"

	"gorm.io/gorm"
)

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password string `json:"password"`
}

func LoginView(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest

		if !isValidData(r.Body, &req) {
			errorResponse(
				&w,
				&errors.AppError{
					Code: 11017,
					Message: "Not valid data",
				},
				http.StatusBadRequest,
		  )
			return
		}

		var user *models.User
		userService := models.UserService{Db: db}

		if !strings.Contains(req.UsernameOrEmail, "@") {
			user = userService.GetByUsername(req.UsernameOrEmail)
		} else {
			user = userService.Objects().Get(func(u *models.User) bool {return u.Email == req.UsernameOrEmail})
		}

		if user == nil {
		  errorResponse(
				&w,
				&errors.AppError{
					Code: 11018,
					Message: "Username/Email is incorrect",
				},
			  http.StatusNotFound,
			)
			return
		}

		if !models.CheckPassword(user.Password, req.Password) {
			errorResponse(
				&w,
				&errors.AppError{
					Code: 11020,
					Message: "Username/Email or password is incorrect ",
				},
				http.StatusUnauthorized,
			)
			return
		}


		if !user.IsVerified {
			response(
				&w,
				map[string]any{
					"success": true,
					"state": 0,
				},
				http.StatusOK,
			)
			return
		}

		token, err := models.GenerateToken(user.Username)

		if err != nil {
			errorResponse(
				&w,
				err,
				http.StatusInternalServerError,
			)
			return
		}

		response(
			&w,
			map[string]any{
				"success": true,
				"state": 0,
				"token": token,
			},
			http.StatusOK,
		)
	}
}