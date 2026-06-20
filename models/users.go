package models

import (
	"fmt"
	"rakkiz-backend/config"
	"rakkiz-backend/errors"
	"rakkiz-backend/validating"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	Username   string    `json:"username" gorm:"primaryKey"`
	Role       RoleType  `json:"role" gorm:"default:user;not null;type:varchar(20);check:role IN ('user','admin','super_admin')"`
	Email      string    `json:"email" gorm:"unique;not null"`
	Name		   string    `json:"name" gorm:"not null"`
	Password   string    `json:"password" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"not null"`
	LastLogin  time.Time `json:"last_login" gorm:"not null"`
	IsVerified bool      `json:"verified" gorm:"not null"`
	Otp        int			 `json:"otp" gorm:"not null"`
	IsMuslim   *bool     `json:"is_muslim" gorm:"not null"`
}

type UserService struct {
  Db *gorm.DB
}

type List[T any] struct {
	All *[]*T
}

func (s *UserService) Create(user *User) (*User, *errors.AppError) {
	if user == nil {
	return nil, &errors.AppError{
			Code:    11019,
			Message: "Missing user data",
		}
	}

	for _, i := range []struct{
		field string
		code  int
		user  *User
	} {
		{
			field: "Username",
			code: 11014,
			user: s.GetByUsername(user.Username),
		},
		{
			field: "Email",
			code: 11015,
			user: s.Objects().Get(func(u *User) bool {return user.Email == u.Email}),
		},
	} {
		if i.user != nil {
			if i.user.IsVerified {
				return nil, &errors.AppError{
					Code:    i.code,
					Message: fmt.Sprintf("%s is already taken", i.field),
				}
			}

			err := i.user.Delete(s.Db)
			if err != nil {return nil, err}
		}
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

	if user.IsMuslim == nil {
		tmp := true
		user.IsMuslim = &tmp
	}

	now := time.Now()

	user.CreatedAt = now
	user.LastLogin = now

	err = errors.Errorize(s.Db.Create(user).Error)
	if err != nil {return nil, err}
	return user, nil
}

func (s *UserService) Objects() *List[User] {
	var users []*User
	s.Db.Find(&users)
	return &List[User]{All: &users}
}

func (s *UserService) GetByUsername(username string) *User {
	var user User

	result := s.Db.Where("username = ?", username).First(&user)

	if result.Error != nil {
		return nil
	}

	return &user
}

func (l *List[T]) Get(fn func(*T) bool) *T {
	for _, item := range *l.All {
		if fn(item) {
			return item
		}
	}
	return nil
}

func (user User) Delete(db *gorm.DB) *errors.AppError {
	return errors.Errorize(db.Unscoped().Delete(user).Error)
}

func (user User) Save(db *gorm.DB) *errors.AppError {
	userService := UserService{Db: db}
	
	usered := userService.GetByUsername(user.Username)

	if usered == nil {
		return &errors.AppError{
			Code:    11016,
			Message: "User not found or username changed",
		}	
	}

	usered = userService.Objects().Get(
		func(u *User) bool {
			return u.Email == user.Email
		},
	)

	if usered != nil && usered.Username != user.Username {
		if usered.IsVerified {
			return &errors.AppError{
				Code:    11015,
				Message: "Email is already taken",
			}
		}
		usered.Delete(db)
	}

	err := ValidateUser(&user)

	if err != nil {return err}

	return errors.Errorize(db.Save(user).Error)
}

func (user User) ChangePasswordAndSave(db *gorm.DB, password string) *errors.AppError {
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
	for _, err := range ([]*errors.AppError{
		validating.ValidateUsername(user.Username),
		validating.ValidateEmail(user.Email),
		validating.ValidateName(user.Name),
		validating.ValidatePassword(user.Password),
	}) {
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateToken(username string) (string, *errors.AppError) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp": time.Now().
				Add(2 * 30 * 24 * time.Hour).
				Unix(),
		},
	)

	t, err := token.SignedString(
		[]byte(config.JWTSecret),
	)
	return t, errors.Errorize(err)
}

func Login(db *gorm.DB, tokenString string) *User {

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (any, error) {
			return []byte(config.JWTSecret), nil
		},
	)

	if err == nil {
		return nil
	}

	userService := UserService{Db: db}
	return userService.GetByUsername(token.Claims.(jwt.MapClaims)["username"].(string))
}