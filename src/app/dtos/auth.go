package dtos

import (
	"time"

	"github.com/pablor21/goms/app/models"
	"github.com/pablor21/goms/pkg/dtos"
	"github.com/pablor21/goms/pkg/interactions/request"
	"github.com/pablor21/goms/pkg/interactions/response"
)

type ClientType string // @name ClientType

const (
	ClientTypeFrontend ClientType = "frontend"
	ClientTypeAdmin    ClientType = "admin"
)

type UserDTO struct {
	dtos.BaseTimestampDTO[int64]
	Email       string            `json:"email,omitempty"`
	PhoneNumber string            `json:"phoneNumber,omitempty"`
	FirstName   string            `json:"firstName,omitempty"`
	LastName    string            `json:"lastName,omitempty"`
	Lang        string            `json:"lang,omitempty"`
	Role        models.UserRole   `json:"role,omitempty"`
	Avatar      *AssetDTO         `json:"avatar,omitempty"`
	Status      models.UserStatus `json:"status,omitempty"`
} // @name User

func (u UserDTO) DisplayName() string {
	if u.FirstName != "" && u.LastName != "" {
		return u.FirstName + " " + u.LastName
	} else if u.FirstName != "" {
		return u.FirstName
	} else if u.LastName != "" {
		return u.LastName
	} else {
		return u.Email
	}
}

type UserResponse struct {
	response.TypedResponse[UserDTO]
} // @name UserResponse

type UserLoginInput struct {
	request.Input
	Username string                 `json:"username" form:"username"`
	Password string                 `json:"password" form:"password"`
	Params   map[string]interface{} `json:"params" form:"params"`
} // @name UserLoginInput

type UserLoginResponseData struct {
	User UserDTO `json:"user"`
} // @name UserResponseLoginData

type UserLoginResponse struct {
	response.TypedResponse[UserLoginResponseData]
} // @name UserLoginResponse

type UpdateProfileInput struct {
	request.Input
	Metadata    dtos.MetadataDTO  `json:"metadata" form:"metadata"`
	FirstName   string            `json:"firstName" form:"firstName"`
	LastName    string            `json:"lastName" form:"lastName"`
	PhoneNumber string            `json:"phoneNumber" form:"phoneNumber"`
	Lang        string            `json:"lang" form:"lang"`
	Avatar      request.FileInput `json:"avatar" form:"avatar"`
} // @name UpdateProfileInput

type UpdateProfileResponse struct {
	response.TypedResponse[UserDTO]
} // @name UpdateProfileResponse

type ChangePasswordInput struct {
	request.Input
	Token              string `json:"token"`
	NewPassword        string `json:"newPassword"`
	NewPasswordConfirm string `json:"newPasswordConfirm"`
} // @name ChangePasswordInput

type ChangePasswordResponse struct {
	response.TypedResponse[bool]
} // @name ChangePasswordResponse

type RequestOTPInput struct {
	request.Input
	Client   ClientType             `json:"client" form:"client"`
	Username string                 `json:"username" form:"username"`
	Params   map[string]interface{} `json:"params" form:"params"`
} // @name RequestOTPInput

type RequestOTPResponseData struct {
	Username    string `json:"username"`
	ResendAfter int64  `json:"resendAfter"`
} // @name RequestOTPResponseData

type RequestOTPResponse struct {
	response.TypedResponse[RequestOTPResponseData]
} // @name RequestOTPResponse

type LoginWithOTPInput struct {
	request.Input
	Username string `json:"username" form:"username" query:"username"`
	Code     string `json:"code" form:"code" query:"code"`
	Token    string `json:"token" form:"token" query:"token"`
} // @name LoginWithOTPInput

// OTPRequestMailData is the data structure used to send the OTP code to the user via email
type OTPResultMailData struct {
	User              UserDTO   `json:"user"`
	Username          string    `json:"username"`
	Expiration        time.Time `json:"expiration"`
	ExpirationMinutes int       `json:"expirationMinutes"`
	RetryIn           int       `json:"retryIn"`
	Token             string    `json:"token"`
	Code              string    `json:"code"`
	Link              string    `json:"link"`
}
