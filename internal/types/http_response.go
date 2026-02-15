package types

// Response represents a standard API response envelope.
type Response[T any] struct {
	Data T      `json:"data,omitempty"`
	Err  string `json:"err,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

// NewSuccessResponse creates a new success response.
func NewSuccessResponse[T any](msg string, data T) Response[T] {
	return Response[T]{
		Msg:  msg,
		Data: data,
	}
}

// NewErrorResponse creates a new error response.
func NewErrorResponse[T any](err string) Response[T] {
	return Response[T]{
		Err: err,
	}
}
