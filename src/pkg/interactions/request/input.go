package request

import (
	"mime/multipart"
)

type Input struct {
} // @name Input

type FileInput struct {
	*multipart.FileHeader `swaggerignore:"true"`
} // @name FileInput

func NewFileInputFromMultipart(f *multipart.FileHeader) *FileInput {
	return &FileInput{
		FileHeader: f,
	}
}

type PaginatedInput struct {
	Input
	Take       int    `json:"take" form:"take" query:"take"`
	Skip       int    `json:"skip" form:"skip" query:"skip"`
	SearchTerm string `json:"searchTerm" form:"searchTerm" query:"searchTerm"`
} // @name PaginatedInput
