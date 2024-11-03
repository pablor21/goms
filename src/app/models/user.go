package models

import (
	"github.com/creasty/defaults"
	"github.com/pablor21/goms/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type UserRole string // @name UserRole

const (
	RoleAdmin UserRole = "ADMIN" //@name Admin
	RoleUser  UserRole = "USER"  //@name User
)

type UserStatus string // @name UserStatus

const (
	UserStatusActive   UserStatus = "ACTIVE"   //@name Active
	UserStatusInactive UserStatus = "INACTIVE" //@name Inactive
)

type User struct {
	models.BaseTimestampedModel[int64]
	models.MetadataModel
	Email         string     `json:"email" gorm:"column:email"`
	PhoneNumber   string     `json:"phone_number" gorm:"column:phone_number"`
	Password      string     `json:"password" gorm:"column:password"`
	FirstName     string     `json:"first_name" gorm:"column:first_name"`
	LastName      string     `json:"last_name" gorm:"column:last_name"`
	Lang          string     `json:"lang" gorm:"column:lang" default:"en"`
	Role          UserRole   `json:"role" gorm:"column:role" default:"USER"`
	SuperAdmin    bool       `json:"super_admin" gorm:"column:super_admin" default:"false"`
	Status        UserStatus `json:"status" gorm:"column:status" default:"ACTIVE"`
	AvatarAssetID *int64     `json:"avatar_asset_id" gorm:"column:avatar_asset_id"`
	Avatar        *Asset     `json:"avatar_asset" gorm:"foreignKey:AvatarAssetID;references:ID"`
	//// Avatar        *Asset     `json:"avatar_asset" gorm:"polymorphicType:OwnerType;polymorphicId:OwnerID;polymorphicValue:users"`
} // @name User

func NewUser() *User {
	u := &User{}
	defaults.Set(u)
	return u
}

func (User) TableName() string {
	return "users"
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) DisplayName() string {
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
