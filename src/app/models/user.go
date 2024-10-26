package models

import (
	"github.com/creasty/defaults"
)

type UserRole string // @name UserRole

const (
	RoleAdmin UserRole = "ADMIN"
	RoleUser  UserRole = "USER"
)

type UserStatus string // @name UserStatus

const (
	StatusActive   UserStatus = "ACTIVE"
	StatusInactive UserStatus = "INACTIVE"
)

type User struct {
	BaseTimestampModel[int64]
	Email       *string    `json:"email" gorm:"column:email"`
	PhoneNumber *string    `json:"phone_number" gorm:"column:phone_number"`
	Password    string     `json:"password" gorm:"column:password"`
	FirstName   *string    `json:"first_name" gorm:"column:first_name"`
	LastName    *string    `json:"last_name" gorm:"column:last_name"`
	Lang        *string    `json:"lang" gorm:"column:lang" default:"en"`
	Role        UserRole   `json:"role" gorm:"column:role" default:"USER"`
	SuperAdmin  bool       `json:"super_admin" gorm:"column:super_admin" default:"false"`
	Status      UserStatus `json:"status" gorm:"column:status" default:"ACTIVE"`
} // @name User

func NewUser() *User {
	u := &User{}
	defaults.Set(u)
	return u
}
