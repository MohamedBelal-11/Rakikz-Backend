package errors

type AppError struct {
	Code    int
	Message string
}

func Errorize(err error) *AppError {
	if err == nil {
		return nil
	}
	return &AppError{
		Code: 12000,
		Message: err.Error(),
	}
}