package interfaceUseCase

import (
	"project/delivery/models"
	"project/domain/entity"
)

type UserUseCase interface {
	CheckAddress(int, string) error
	ExecuteAddAddress(*entity.UserAddress) error
	ExecuteChangePassword(int) (string, error)
	ExecuteDeleteAddress(int, string) error
	ExecuteEditAddress(entity.UserAddress, int, string) error
	ExecuteEditProfile(entity.User, int) error
	ExecuteLogin(string) (string, error)
	ExecuteLoginWithPassword(string, string) (int, error)
	ExecuteOtpValidation(string, string) (*entity.User, error)
	ExecuteOtpValidationPassword(string, string, int) error
	ExecuteShowUserDetails(int) (*entity.User, *entity.UserAddress, error)
	ExecuteSignup(entity.User) (*entity.User, error)
	ExecuteSignupOtpValidation( string,  string)  (int,string,error)
	ExecuteSignupWithOtp(models.Signup) (string, error)
	ExecuteForgetPassword( string) (string, error)
	ExecuteOtpValidationFPassword( string,  string,  string) error
}
