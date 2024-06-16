package interfaces

import (
	"project/delivery/models"
	"project/domain/entity"
)

type UserRepository interface {
	CheckPermission(*entity.User) (bool, error)
	Create(*entity.User) error
	CreateAddress(*entity.UserAddress) error
	CreateOtpKey(string, string) error
	CreateSignup(*models.Signup) error
	Delete(*entity.User) error
	DeleteAddress(int) error
	GetAddressByID(int) (*entity.UserAddress, error)
	GetAddressById(int) (*entity.UserAddress, error)
	GetAddressByType(int, string) (*entity.UserAddress, error)
	GetByEmail(string) (*entity.User, error)
	GetById(int) (*entity.User, error)
	GetByKey(string) (*entity.OtpKey, error)
	GetByPhone(string) (*entity.User, error)
	GetByReferalCode(string) (*entity.User, error)
	GetSignupByPhone(string) (*models.Signup, error)
	Update(*entity.User) error
	UpdateAddress(*entity.UserAddress) error
	UpdateOtp( *entity.OtpKey) error
	CheckValidation( string) (bool,error)
	PhoneExists( string) (bool, error)
}
