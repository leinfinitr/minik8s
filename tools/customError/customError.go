package customError

import (
	"fmt"
)

type CustomError struct {
	Message string
}

// 实现error接口的Error方法
func (e *CustomError) Error() string {
	return fmt.Sprintf(e.Message)
}
