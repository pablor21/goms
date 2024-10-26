package dtos

type BasDTO[IDType any] struct {
	ID IDType `json:"id"`
}

type BaseTimestampDTO[IDType any] struct {
	BasDTO[IDType]
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type Input struct {
}

type Response[T any] struct {
	Result        T      `json:"result,omitempty"`
	StatusCode    int    `json:"status_code,omitempty" default:"200"`
	StatusMessage string `json:"status_message,omitempty" default:"Ok"`
	Message       string `json:"message"`
	Error         error  `json:"error"`
} // @name Response[T]
