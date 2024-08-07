package usecase

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"regexp"

	"project/config"
	"project/delivery/models"
	"project/domain/entity"
	"project/domain/utils"
	interfaces "project/repository/interfaceRepository"
	"project/usecase/interfaceUseCase"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo interfaces.UserRepository
	otp      *config.OTP
}

func NewUser(userRepo interfaces.UserRepository, otp *config.OTP) interfaceUseCase.UserUseCase {
	return &UserUseCase{userRepo: userRepo, otp: otp}
}

func (us *UserUseCase) ExecuteSignup(user entity.User) (*entity.User, error) {

	email, err := us.userRepo.GetByEmail(user.Email)
	if err != nil {
		return nil, errors.New("error with server")
	}
	if email != nil {
		return nil, errors.New("user with email already exists")
	}
	phone, err := us.userRepo.GetByPhone(user.Phone)
	if err != nil {
		return nil, errors.New("error eith server")
	}
	if phone != nil {
		return nil, errors.New("user with phoone already exists")
	}
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &entity.User{
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Password: string(hashedpassword),
	}

	err1 := us.userRepo.Create(newUser)
	if err1 != nil {
		return nil, errors.New("Error creating user")
	}
	return newUser, nil
}
func isValidName(name string) bool {
	alphaRegex := regexp.MustCompile("^[a-zA-Z]+$")
	return alphaRegex.MatchString(name)
}

func (uu *UserUseCase) ExecuteSignupWithOtp(user models.Signup) (string, error) {
	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return "", err
		}
		errors := err.(validator.ValidationErrors)
		errorMsg := "Validation failed: "
		for _, e := range errors {
			switch e.Tag() {
			case "required":
				errorMsg += fmt.Sprintf("%s is required; ", e.Field())
			case "alpha":
				errorMsg += fmt.Sprintf("%s should contain only alphabetic characters; ", e.Field())
			case "email":
				errorMsg += fmt.Sprintf("%s should be a valid email; ", e.Field())
			case "numeric":
				errorMsg += fmt.Sprintf("%s should contain only numeric values; ", e.Field())
			case "min":
				errorMsg += fmt.Sprintf("%s should contain minimum 8 charchters; ", e.Field())
			case "len":
				errorMsg += fmt.Sprintf("%s should contain 10 numbers; ", e.Field())
			default:
				errorMsg += fmt.Sprintf("%s has an invalid value; ", e.Field())
			}
		}
		return "", fmt.Errorf(errorMsg)
	}
	var otpKey entity.OtpKey

	phone, err := uu.userRepo.GetByPhone(user.Phone)
	if err != nil {
		return "", errors.New("error with server")
	}
	if phone != nil {
		return "", errors.New("user with this phone no already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.Password = string(hashedPassword)
	if user.ReferalCode != "" {
		_, err := uu.userRepo.GetByReferalCode(user.ReferalCode)
		if err != nil {
			return "", errors.New("Invalid referal code")
		}
	}

	key, err := utils.SendOtp(user.Phone, *uu.otp)
	if err != nil {
		return "", err
	} else {
		err = uu.userRepo.CreateSignup(&user)
		otpKey.Key = key
		otpKey.Phone = user.Phone
		err = uu.userRepo.CreateOtpKey(key, user.Phone)
		if err != nil {
			return "", err
		}
		return key, nil
	}
}

func (uu *UserUseCase) ExecuteSignupOtpValidation(key string, otp string) (int, string, error) {
	result, err := uu.userRepo.GetByKey(key)

	if err != nil {
		return 0, "", errors.New("error in key")
	}
	fmt.Printf("GetByKey Result: %+v\n", result)
	user, err := uu.userRepo.GetSignupByPhone(result.Phone)
	if err != nil {
		return 0, "", errors.New("error in phone")
	}
	err = utils.CheckOtp(result.Phone, otp, *uu.otp)
	if err != nil {
		return 0, "", err
	} else {
		newUser := &entity.User{
			Name:     user.Name,
			Email:    user.Email,
			Phone:    user.Phone,
			Password: user.Password,
		}
		err1 := uu.userRepo.Create(newUser)
		if err1 != nil {
			return 0, "", errors.New("error while crearting user")
		}
		if user.ReferalCode != "" {
			referduser, err := uu.userRepo.GetByReferalCode(user.ReferalCode)
			if err != nil {
				return 0, "", err
			}
			if referduser != nil {
				referduser.Wallet = referduser.Wallet + 500
				err = uu.userRepo.Update(referduser)
				newUser.Wallet = 500
			}
		}
		randomBytes := make([]byte, 3)
		_, err2 := rand.Read(randomBytes)
		if err2 != nil {
			return 0, "", err
		}

		referralCode := hex.EncodeToString(randomBytes)

		newUser.ReferalCode = referralCode

		err3 := uu.userRepo.Update(newUser)
		if err3 != nil {
			return 0, "", errors.New("user update failed")
		}

		return newUser.Id, newUser.Phone, nil
	}

}

func (uu *UserUseCase) ExecuteLoginWithPassword(phone, password string) (int, error) {

	user, err := uu.userRepo.GetByPhone(phone)
	if err != nil {
		return 0, err
	}
	if user == nil {
		return 0, errors.New("user with this phone not found")
	}

	permission, err := uu.userRepo.CheckPermission(user)
	if err != nil {
		return 0, err
	}
	if permission == false {
		return 0, errors.New("permission denied")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return 0, errors.New("invalid password")
	} else {
		return user.Id, nil
	}
}

func (u *UserUseCase) ExecuteLogin(phone string) (string, error) {

	var otpKey entity.OtpKey
	result, err := u.userRepo.GetByPhone(phone)
	if err != nil {
		return "", err
	}
	if result == nil {
		return "", errors.New("no user with this phone found")
	}
	permission, err := u.userRepo.CheckPermission(result)
	if err != nil {
		return "", err
	}
	if permission == false {
		return "", errors.New("permission denied")
	}
	key, err := utils.SendOtp(phone, *u.otp)
	if err != nil {
		return "", err
	} else {
		otpKey.Key = key
		otpKey.Phone = phone
		err = u.userRepo.CreateOtpKey(key, phone)
		if err != nil {
			return "", err
		}
		return key, nil
	}

}

func (uu *UserUseCase) ExecuteOtpValidation(key, otp string) (*entity.User, error) {
	result, err := uu.userRepo.GetByKey(key)
	if err != nil {
		return nil, err
	}
	user, err := uu.userRepo.GetByPhone(result.Phone)
	if err != nil {
		return nil, err
	}
	err1 := utils.CheckOtp(result.Phone, otp, *uu.otp)
	if err1 != nil {
		return nil, err
	}
	return user, nil
}

func (uu *UserUseCase) ExecuteAddAddress(address *entity.UserAddress) error {
	validate := validator.New()
	if err := validate.Struct(address); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		errors := err.(validator.ValidationErrors)
		errorMsg := "Validation failed: "
		for _, e := range errors {
			switch e.Tag() {
			case "required":
				errorMsg += fmt.Sprintf("%s is required; ", e.Field())
			case "alpha":
				errorMsg += fmt.Sprintf("%s should contain only alphabetic characters; ", e.Field())
			default:
				errorMsg += fmt.Sprintf("%s has an invalid value; ", e.Field())
			}
		}
		return fmt.Errorf(errorMsg)
	}

	err := uu.userRepo.CreateAddress(address)
	if err != nil {
		return err
	}
	return nil

}

func (uu *UserUseCase) CheckAddress(id int, adtype string) error {
	ad, err := uu.userRepo.GetAddressByType(id, adtype)
	if err != nil {
		return errors.New("err geting address")
	}
	if ad != nil {
		return errors.New("address with the type already exists")
	}
	return nil
}

func (uu *UserUseCase) ExecuteEditProfile(user entity.User, userid int) error {
	// validate := validator.New()
	// if err := validate.Struct(user); err != nil {
	// 	if _, ok := err.(*validator.InvalidValidationError); ok {
	// 		return err
	// 	}
	// 	errors := err.(validator.ValidationErrors)
	// 	errorMsg := "Validation failed: "
	// 	for _, e := range errors {
	// 		switch e.Tag() {
	// 		case "required":
	// 			errorMsg += fmt.Sprintf("%s is required; ", e.Field())
	// 		case "alpha":
	// 			errorMsg += fmt.Sprintf("%s should contain only alphabetic characters; ", e.Field())
	// 		case "email":
	// 			errorMsg += fmt.Sprintf("%s should be valid email; ", e.Field())
	// 		default:
	// 			errorMsg += fmt.Sprintf("%s has an invalid value; ", e.Field())
	// 		}
	// 	}
	// 	return fmt.Errorf(errorMsg)
	// }
	user.Id = userid
	err := uu.userRepo.Update(&user)
	if err != nil {
		return errors.New("useer updation failed")
	}
	return nil
}

func (uu *UserUseCase) ExecuteShowUserDetails(userid int) (*entity.User, *entity.UserAddress, error) {
	user, err := uu.userRepo.GetById(userid)
	if err != nil {
		return nil, nil, err
	}
	address, err1 := uu.userRepo.GetAddressByID(userid)
	if err1 != nil {
		return nil, nil, err1
	}
	if user != nil && address != nil {
		return user, address, nil
	} else {
		return nil, nil, errors.New("user with this id not found")
	}
}

func (uu *UserUseCase) ExecuteChangePassword(userid int) (string, error) {
	var otpkey entity.OtpKey
	user, err := uu.userRepo.GetById(userid)
	if err != nil {
		return "", err
	}
	key, err1 := utils.SendOtp(user.Phone, *uu.otp)
	if err1 != nil {
		return "", err1
	} else {
		otpkey.Key = key
		otpkey.Phone = user.Phone
		err := uu.userRepo.CreateOtpKey(otpkey.Key, otpkey.Phone)
		if err != nil {
			return "", nil
		}
		return key, nil
	}
}

func (uu *UserUseCase) ExecuteForgetPassword(phone string) (string, error) {
	var otpkey entity.OtpKey
	user, err := uu.userRepo.PhoneExists(phone)
	if err != nil {
		return "", err
	}
	if user == true {
		key, err1 := utils.SendOtp(phone, *uu.otp)
		if err1 != nil {
			return "", err1
		}
		otpkey.Key = key
		otpkey.Phone = phone
		err = uu.userRepo.CreateOtpKey(otpkey.Key, otpkey.Phone)
		if err != nil {
			return "", nil
		}
		return key, nil
	}
	return "", errors.New("phone doesnt exists")
}

func (uu *UserUseCase) ExecuteOtpValidationPassword(password string, otp string, userid int) error {
	user, err := uu.userRepo.GetById(userid)
	if err != nil {
		return err
	}
	err = utils.CheckOtp(user.Phone, otp, *uu.otp)
	if err != nil {
		return err
	}
	hashedpassword, err1 := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.Password = string(hashedpassword)
	err1 = uu.userRepo.Update(user)
	if err1 != nil {
		return errors.New("password changing failed")

	}
	return nil

}

func (uu *UserUseCase) ExecuteOtpValidationFPassword(otp string, key string) error {
	phone, err := uu.userRepo.GetByKey(key)
	if err != nil {
		return err
	}
	// user, err := uu.userRepo.GetByPhone(phone.Phone)
	// if err != nil {
	// 	return err
	// }
	err = utils.CheckOtp(phone.Phone, otp, *uu.otp)
	if err != nil {
		return err
	}
	phone.Validated = true
	err = uu.userRepo.UpdateOtp(phone)
	if err != nil {
		return err
	}
	fmt.Println("ressss", phone.Validated)

	// hashedpassword, err1 := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// user.Password = string(hashedpassword)
	// err1 = uu.userRepo.Update(user)
	// if err1 != nil {
	// 	return errors.New("password changing failed")

	// }
	return nil

}

func (uu *UserUseCase) ForgetPassChange(key string, password string) error {
	res, err := uu.userRepo.CheckValidation(key)
	if err != nil {
		return err
	}
	fmt.Println("res", res)
	phone, err := uu.userRepo.GetByKey(key)
	if err != nil {
		return err
	}
	fmt.Println("ph", phone.Phone)
	user, err := uu.userRepo.GetByPhone(phone.Phone)
	if err != nil {
		return err
	}

	if res == true {
		hashedpassword, err1 := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		user.Password = string(hashedpassword)
		err1 = uu.userRepo.Update(user)
		if err1 != nil {
			return errors.New("password changing failed")

		}
	} else {
		return errors.New("Otp Not validated")
	}
	return nil
}

func (uu *UserUseCase) ExecuteEditAddress(usaddress entity.UserAddress, id int, useraddress string) error {
	validate := validator.New()
	if err := validate.Struct(usaddress); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		errors := err.(validator.ValidationErrors)
		errorMsg := "Validation failed: "
		for _, e := range errors {
			switch e.Tag() {
			case "required":
				errorMsg += fmt.Sprintf("%s is required; ", e.Field())
			case "alpha":
				errorMsg += fmt.Sprintf("%s should contain only alphabetic characters; ", e.Field())
			case "len":
				errorMsg += fmt.Sprintf("%s should have a length of %s; ", e.Field(), e.Param())
			case "numeric":
				errorMsg += fmt.Sprintf("%s should contain only numeric characters; ", e.Field())
			default:
				errorMsg += fmt.Sprintf("%s has an invalid value; ", e.Field())
			}
		}
		return fmt.Errorf(errorMsg)
	}
	exisitingaddress, err := uu.userRepo.GetAddressByType(id, useraddress)
	if err != nil {
		return err
	}

	exisitingaddress.User_id = id
	exisitingaddress.Address = usaddress.Address
	exisitingaddress.State = usaddress.State
	exisitingaddress.Country = usaddress.Country
	exisitingaddress.Pin = usaddress.Pin
	exisitingaddress.Type = usaddress.Type

	err1 := uu.userRepo.UpdateAddress(exisitingaddress)
	if err1 != nil {
		return err1
	}
	return nil
}

func (uu *UserUseCase) ExecuteDeleteAddress(id int, addtype string) error {
	address, err := uu.userRepo.GetAddressByType(id, addtype)
	if err != nil {
		return err
	}
	err1 := uu.userRepo.DeleteAddress(address.Id)
	if err1 != nil {
		return err1
	}
	return nil
}

func (uu *UserUseCase) GetUserAddressByID(id int) (entity.UserAddress, error) {
	addres, err := uu.userRepo.GetAddressById(id)
	if err != nil {
		return entity.UserAddress{}, err
	}
	return *addres, nil
}

func (uu *UserUseCase) GetAllUserAddress(id int) ([]entity.UserAddress, error) {
	add, err := uu.userRepo.GetAddressesByUserID(id)
	if err != nil {
		return nil, err
	}
	return add, nil
}
