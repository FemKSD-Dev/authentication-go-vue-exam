package response

type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Data    *T     `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

func Success[T any](data T, message string) *APIResponse[T] {
	if message == "" {
		message = "Operation completed successfully"
	}
	return &APIResponse[T]{
		Success: true,
		Code:    "OK",
		Data:    &data,
		Message: message,
	}
}

func Fail[T any](code, message string) *APIResponse[T] {
	if message == "" {
		message = "An error occurred"
	}
	if code == "" {
		code = CodeInternalError
	}
	return &APIResponse[T]{
		Success: false,
		Code:    code,
		Message: message,
	}
}
