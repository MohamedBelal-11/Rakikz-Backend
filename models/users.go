package models

import (
	"rakkiz-backend/bstrings"
	"rakkiz-backend/errors"
	"rakkiz-backend/validating"
	"slices"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)
type RoleType string

const (
	RoleUser RoleType = "user"
	RoleAdmin RoleType = "admin"
	RoleSuperAdmin RoleType = "super_admin"
)


type User struct {
	Username   string    `json:"username" gorm:"primaryKey;index"`
	Role       RoleType  `json:"role" gorm:"default:user; not null;type:varchar(20);check:role IN ('user','admin','super_admin')"`
	Email      string    `json:"email" gorm:"unique;not null"`
	Name		   string    `json:"name" gorm:"not null"`
	Password   string    `json:"password" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"not null"`
	LastLogin  time.Time `json:"last_login" gorm:"not null"`
	IsVerified bool      `json:"verified" gorm:"not null"`
	Otp        int			 `json:"otp" gorm:"not null"`
}

type UserService struct {
  Db *gorm.DB
}

type List[T any] struct {
	All []T
}

func (s *UserService) Create(user *User) (*User, *errors.AppError) {
	usered := s.GetByUsername(user.Username)

	if usered != nil {
		if usered.IsVerified {
			return nil, &errors.AppError{
				Code:    11014,
				Message: "Username is already taken",
			}
		}
		err := usered.Delete(s.Db)
		if err != nil {return nil, err}
	}

	if usered != nil {
		if usered.IsVerified {
			return nil, &errors.AppError{
				Code:    11015,
				Message: "Email is already taken",
			}
		}
		err := usered.Delete(s.Db)
		if err != nil {return nil, err}
	}

	err := ValidateUser(user)

	if err != nil {return nil, err}

	user.Password, err = HashPassword(user.Password)

	if err != nil {return nil, err}

	switch user.Role {
		case RoleUser, RoleAdmin, RoleSuperAdmin:
		default:
			user.Role = RoleUser
	}

	now := time.Now()

	user.CreatedAt = now
	user.LastLogin = now

	err = errors.Errorize(s.Db.Create(user).Error)
	if err != nil {return nil, err}
	return user, nil
}

func (s *UserService) Objects() List[User] {
	var users []User
	s.Db.Find(&users)
	return List[User]{All: users}
}

func (s *UserService) GetByUsername(username string) *User {
	var user User

	result := s.Db.Where("username = ?", username).First(&user)

	if result.Error != nil {
		return nil
	}

	return &user
}

func (l *List[T]) Get(fn func(T) bool) *T {
	for _, item := range l.All {
		if fn(item) {
			return &item
		}
	}
	return nil
}

func (user *User) Delete(db *gorm.DB) *errors.AppError {
	return errors.Errorize(db.Unscoped().Delete(user).Error)
}

func (user *User) Save(db *gorm.DB) *errors.AppError {
	userService := UserService{Db: db}
	users := userService.Objects()
	usered := users.Get(func(u User) bool { return u.Username == user.Username })

	if usered == nil {
		return &errors.AppError{
			Code:    11016,
			Message: "User not found or username changed",
		}	
	}

	usered = users.Get(func(u User) bool { return u.Email == user.Email })
	if usered != nil && usered.Username != user.Username {
		if usered.IsVerified {
			return &errors.AppError{
				Code:    11015,
				Message: "Email is already taken",
			}
		}
		usered.Delete(db)
	}

	err := ValidateUser(user)

	if err != nil {return err}

	return errors.Errorize(db.Save(user).Error)
}

func (user *User) ChangePasswordAndSave(db *gorm.DB, password string) *errors.AppError {
	pass, err := HashPassword(password)
	if err != nil {return err}

	user.Password = pass

	return user.Save(db)
}

func HashPassword(password string) (string, *errors.AppError) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(hash), errors.Errorize(err)
}

func CheckPassword(hash, password string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func ValidateUser(user *User) *errors.AppError {
	if user.Username == "" {
		return &errors.AppError{
			Code:    11000,
			Message: "Username is required",
		}
	}

	if len(user.Username) < 3 || len(user.Username) > 25 {
		return &errors.AppError{
			Code:    11006,
			Message: "Username must be between 3 and 25 characters",
		}
	}

	allowed, char := validating.IsFromChars(user.Username, usernameAllowedChars)
	if allowed == false {
		return &errors.AppError{
			Code:    11001,
			Message: "Username contains invalid characters: " + char,
		}
	}

	if user.Email == "" {
		return &errors.AppError{
			Code:    11002,
			Message: "Email is required",
		}
	}

	if len(user.Email) > 100 {
		return &errors.AppError{
			Code:    11008,
			Message: "Email must be less than 100 characters",
		}
	}

	if validating.HasNotAllowedSpace(user.Email) {
		return &errors.AppError{
			Code:    11013,
			Message: "Email contains not allowed spaces",
		}
	}

	if !validating.IsValidEmail(user.Email) {
		return &errors.AppError{
			Code:    11003,
			Message: "Email is invalid",
		}
	}

	if user.Name == "" {
		return &errors.AppError{
			Code:    11004,
			Message: "Name is required",
		}
	}

	if validating.HasNotAllowedSpace(user.Name) {
		return &errors.AppError{
			Code:    11005,
			Message: "Name contains not allowed spaces",
		}
	} 

	if len(user.Name) < 3 || len(user.Name) > 25 {
		return &errors.AppError{
			Code:    11007,
			Message: "Name must be between 3 and 25 characters",
		}
	}
	allowed, char = hasAllowedChars(user.Name)
	if !allowed {
		return &errors.AppError{
			Code:    1109,
			Message: "Name contains invalid characters",
		}
	}

	if user.Password == "" {
		return &errors.AppError{
			Code:    11010,
			Message: "Password is required",
		}
	}

	if len(user.Password) < 8 || len(user.Password) > 200 {
		return &errors.AppError{
			Code:    11011,
			Message: "Password must be between 8 and 200 characters",
		}
	}

	allowed, char = validating.IsFromChars(user.Password, passwordAllowedChars)
	if !allowed {
		return &errors.AppError{
			Code:    11012,
			Message: "Password contains invalid characters: " + char,
		}
	}
	return nil
}


var usernameAllowedChars = slices.Concat(
	bstrings.AllowedUsernameMarks,
	bstrings.EnLittters,
	bstrings.Numbers,
)

var passwordAllowedChars = slices.Concat(
	bstrings.AllowedPasswordMarks,
	bstrings.EnLittters,
	bstrings.Numbers,
)

func hasAllowedChars(name string) (bool, string) {
	for _, ch := range name {
		if !unicode.IsLetter(ch) && !unicode.IsNumber(ch) && !(string(ch) == " ") {
			return false, string(ch)
		}
	}
	return true, ""
}