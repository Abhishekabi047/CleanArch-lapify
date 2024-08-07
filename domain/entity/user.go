package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model  `json:"-"`
	Id          int    `gorm:"primarykey"  json:"id"`
	Name        string `json:"name"  `
	Email       string `json:"email" `
	Phone       string `json:"phone" validate:"required"`
	Password    string `json:"password" validate:"required,min=8"`
	IsBlocked   bool   `gorm:"not null;default:true" json:"-"`
	Wallet      int    `json:"wallet"`
	Permission  bool   `gorm:"not null;default:true" json:"permission"`
	ReferalCode string `json:"referalcode"`
}

type UserAddress struct {
	gorm.Model `json:"-"`
	Id         int    `gorm:"primarykey" json:"id"`
	User_id    int    `json:"-"`
	FullName   string `json:"full_name"`
	Phone      int    `json:"phone" validate:"required"`
	AltPhone   int    `json:"alt_phone" validate:"required"`
	Address    string `json:"address" validate:"required"`
	State      string `json:"state" validate:"required,alpha"`
	Country    string `json:"country" validate:"required,alpha"`
	Pin        string `json:"pin" validate:"required,numeric,len=6"`
	Type       string `json:"type" validate:"required,alpha"`
}

type Login struct {
	Email    string `json:"email" bson:"email" binding:"required"`
	Password string `json:"pasword" bson:"password" binding:"required"`
}

type OtpKey struct {
	gorm.Model
	Key       string `json:"key"`
	Phone     string `json:"phone"`
	Validated bool   `json:"validated" gorm:"default:false"`
}
type ListUsersResponse struct {
	Users []User `json:"users"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
