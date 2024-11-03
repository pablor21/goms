package response

import (
	"net/http"

	"github.com/pablor21/goms/pkg/errors"
)

type Response struct {
	StatusCode    int                    `json:"statusCode,omitempty" default:"200"`
	StatusMessage string                 `json:"statusMessage,omitempty" default:"Ok"`
	Message       string                 `json:"message,omitempty"`
	Error         error                  `json:"error,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
} // @name Response

func NewResponse() Response {
	return Response{
		StatusCode:    200,
		StatusMessage: "Ok",
	}
}

func (r Response) SetCode(code int) Response {
	r.StatusCode = code
	r.StatusMessage = http.StatusText(code)
	return r
}

func (r Response) SetMessage(message string) Response {
	r.Message = message
	return r
}

func (r Response) SetError(err error) Response {
	r.Error = err
	if err != nil {
		if appError, ok := err.(*errors.AppError); ok {
			r.StatusCode = appError.StatusCode
			r.StatusMessage = appError.StatusMessage
		}
	}
	return r
}

func (r Response) SetMetadata(metadata map[string]interface{}) Response {
	r.Metadata = metadata
	return r
}

func (r Response) SetMetaValue(key string, value interface{}) Response {
	if r.Metadata == nil {
		r.Metadata = make(map[string]interface{})
	}
	r.Metadata[key] = value
	return r
}

func (r Response) GetMetaValue(key string) (interface{}, bool) {
	if r.Metadata == nil {
		return nil, false
	}
	value, ok := r.Metadata[key]
	return value, ok
}

func (r Response) SetStatusMessage(statusMessage string) Response {
	r.StatusMessage = statusMessage
	return r
}

func Ok() Response {
	return NewResponse().SetCode(200)
}

// / TYPED RESPONSE
type TypedResponse[T any] struct {
	Response
	Result T `json:"result,omitempty"`
} // @name TypedResponse[T]

func NewTypedResponse[T any](result T) TypedResponse[T] {
	return TypedResponse[T]{
		Result: result,
	}.SetCode(200)
}

func (r TypedResponse[T]) SetResult(result T) TypedResponse[T] {
	r.Result = result
	return r
}

func (r TypedResponse[T]) SetCode(code int) TypedResponse[T] {
	r.Response = r.Response.SetCode(code)
	return r
}

func (r TypedResponse[T]) SetMessage(message string) TypedResponse[T] {
	r.Response = r.Response.SetMessage(message)
	return r
}

func (r TypedResponse[T]) SetError(err error) TypedResponse[T] {
	r.Response = r.Response.SetError(err)
	return r
}

func (r TypedResponse[T]) SetMetadata(metadata map[string]interface{}) TypedResponse[T] {
	r.Response = r.Response.SetMetadata(metadata)
	return r
}

func (r TypedResponse[T]) SetMetaValue(key string, value interface{}) TypedResponse[T] {
	r.Response = r.Response.SetMetaValue(key, value)
	return r
}

func (r TypedResponse[T]) GetMetaValue(key string) (interface{}, bool) {
	return r.Response.GetMetaValue(key)
}

func (r TypedResponse[T]) SetStatusMessage(statusMessage string) TypedResponse[T] {
	r.Response = r.Response.SetStatusMessage(statusMessage)
	return r
}

func TOk[T any](result T) TypedResponse[T] {
	return NewTypedResponse(result).SetCode(200)
}

type PaginatedResponseData[T any] struct {
	Items   []T  `json:"items,omitempty"`
	Total   int  `json:"total,omitempty"`
	First   int  `json:"first,omitempty"`
	Last    int  `json:"last,omitempty"`
	HasMore bool `json:"hasMore,omitempty"`
} // @name PaginatedResponseData[T]

func NewPaginatedResponseData[T any](items []T, total, first, last int) PaginatedResponseData[T] {
	return PaginatedResponseData[T]{
		Items:   items,
		Total:   total,
		First:   first,
		Last:    last,
		HasMore: last < total,
	}
}

type PaginatedResponse[T any] struct {
	TypedResponse[PaginatedResponseData[T]]
} // @name PaginatedResponse[T]
