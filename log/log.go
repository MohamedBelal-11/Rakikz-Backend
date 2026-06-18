package log

import (
	"fmt"
	"rakkiz-backend/errors"
)

func Erorr(err *errors.AppError) {
	if err == nil {
		fmt.Println("No Error")
	} else {
		fmt.Printf("Error: [%d] %s\n", err.Code, err.Message)
	}
}

