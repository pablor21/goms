package dtos

import "github.com/pablor21/goms/app/models"

type UserDTO struct {
	ID          int               `json:"id"`
	Email       string            `json:"email,omitempty"`
	PhoneNumber string            `json:"phone_number,omitempty"`
	FirstName   string            `json:"first_name,omitempty"`
	LastName    string            `json:"last_name,omitempty"`
	Lang        string            `json:"lang,omitempty"`
	Role        models.UserRole   `json:"role,omitempty"`
	Status      models.UserStatus `json:"status,omitempty"`
	CreatedAt   string            `json:"created_at,omitempty"`
	UpdatedAt   string            `json:"updated_at,omitempty"`
} // @name User

type UserLoginInput struct {
	FieldName  string                 `json:"field_name"`
	FieldValue string                 `json:"field_value"`
	Data       map[string]interface{} `json:"data"`
} // @name UserLoginInput

type UserResponseLoginData struct {
	User UserDTO `json:"user"`
} // @name UserResponseLoginData

type UserLoginResponse struct {
	Response[UserResponseLoginData]
} // @name UserLoginResponse
