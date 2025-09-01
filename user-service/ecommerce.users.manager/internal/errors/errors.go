package errors

type CodeError int

const (
	ErrCodeInvalidArgument CodeError = iota
	ErrCodeNotFound
	ErrCodeInternal
)

type Error struct {
	Code    CodeError
	Message string
}

func (e Error) Error() string {
	return e.Message
}

func New(code CodeError, message string) error {
	return Error{
		Code:    code,
		Message: message,
	}
}