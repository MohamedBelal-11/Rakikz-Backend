package validating

import (
	"net/mail"
	"rakkiz-backend/bstrings"
	"rakkiz-backend/errors"
	"slices"
	"unicode"
)

func ValidateUsername(username string) *errors.AppError {
	if username == "" {
		return &errors.AppError{
			Code:    11000,
			Message: "Username is required",
		}
	}

	if len(username) < 3 || len(username) > 25 {
		return &errors.AppError{
			Code:    11006,
			Message: "Username must be between 3 and 25 characters",
		}
	}

	allowed, char := IsFromChars(username, usernameAllowedChars)
	if allowed == false {
		return &errors.AppError{
			Code:    11001,
			Message: "Username contains invalid characters: " + char,
		}
	}
	return nil
}

func ValidateEmail(email string) *errors.AppError {
	if email == "" {
		return &errors.AppError{
			Code:    11002,
			Message: "Email is required",
		}
	}

	if len(email) > 100 {
		return &errors.AppError{
			Code:    11008,
			Message: "Email must be less than 100 characters",
		}
	}

	if HasNotAllowedSpace(email) {
		return &errors.AppError{
			Code:    11013,
			Message: "Email contains not allowed spaces",
		}
	}

	if !IsValidEmail(email) {
		return &errors.AppError{
			Code:    11003,
			Message: "Email is invalid",
		}
	}
	return nil
}

func ValidateName(name string) *errors.AppError {
	
	if name == "" {
		return &errors.AppError{
			Code:    11004,
			Message: "Name is required",
		}
	}

	if HasNotAllowedSpace(name) {
		return &errors.AppError{
			Code:    11005,
			Message: "Name contains not allowed spaces",
		}
	} 

	if len(name) < 3 || len(name) > 25 {
		return &errors.AppError{
			Code:    11007,
			Message: "Name must be between 3 and 25 characters",
		}
	}
	allowed, char := hasAllowedChars(name)
	if !allowed {
		return &errors.AppError{
			Code:    1109,
			Message: "Name contains invalid characters: " + char,
		}
	}
	return nil
}

func ValidatePassword(password string) *errors.AppError {
	if password == "" {
		return &errors.AppError{
			Code:    11010,
			Message: "Password is required",
		}
	}

	if len(password) < 8 || len(password) > 200 {
		return &errors.AppError{
			Code:    11011,
			Message: "Password must be between 8 and 200 characters",
		}
	}

	allowed, char := IsFromChars(password, passwordAllowedChars)
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
		if !unicode.IsLetter(ch) && !unicode.IsNumber(ch) && string(ch) != " " {
			return false, string(ch)
		}
	}
	return true, ""
}


func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}