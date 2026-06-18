package log

import (
	"fmt"
	"rakkiz-backend/errors"
)

func Erorr(err *errors.AppError) {
	if err == nil {
		fmt.Println("No Error")
	} else {
		fmt.Println(Serror(err))
	}
}

func Serror (err *errors.AppError) string {
	return fmt.Sprintf("Error: [%d] %s", err.Code, err.Message)
}

