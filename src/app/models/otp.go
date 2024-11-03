package models

import (
	"time"

	"github.com/pablor21/goms/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type OTP struct {
	models.BaseTimestampedModel[int64]
	models.MetadataModel
	Username      string     `json:"username" gorm:"column:username"`
	UserID        int64      `json:"user_id" gorm:"column:user_id"`
	User          *User      `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Code          string     `json:"code" gorm:"column:code"`
	MaxAttempts   int        `json:"maxAttempts" gorm:"column:max_attempts"`
	AttemptsCount int        `json:"attemptsCount" gorm:"column:attempts_count"`
	ValidUntil    *time.Time `json:"validUntil" gorm:"column:valid_until"`
}

func (OTP) TableName() string {
	return "otps"
}

func (m *OTP) IsExpired() bool {
	return m.ValidUntil != nil && m.ValidUntil.Before(time.Now())
}

func (u *OTP) CheckCode(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Code), []byte(password))
	return err == nil
}

func (u *OTP) SetCode(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Code = string(hashedPassword)
	return nil
}
