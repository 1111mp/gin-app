package errors

import "fmt"

// RepositoryError - 表示数据库或持久化层的错误
type RepositoryError struct {
	Message string
	Err     error
}

func (e *RepositoryError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s", e.Err)
	}
	return e.Message
}

func (e *RepositoryError) Unwrap() error {
	return e.Err
}

func NewRepositoryError(msg string, err error) *RepositoryError {
	return &RepositoryError{
		Message: msg,
		Err:     err,
	}
}
