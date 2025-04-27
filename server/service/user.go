package service

import (
	"errors"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"server/global"
	"server/model/apptypes"
	"server/model/database"
	"server/model/request"
	"server/model/response"
	"server/utils"
)

type UserService struct {
}

func (userService *UserService) Register(req request.Register) (database.User, error) {
	var user database.User
	if !errors.Is(global.DB.Where("email = ?", req.Email).First(&user).Error, gorm.ErrRecordNotFound) {
		return database.User{}, errors.New("this email address is already registered")
	}
	user.Email = req.Email
	user.Password = utils.BcryptHash(req.Password)
	user.UUID = uuid.Must(uuid.NewV4())
	user.Avatar = "/image/avatar.jpg"
	user.RoleID = apptypes.User
	user.Register = apptypes.Email

	if err := global.DB.Create(&user).Error; err != nil {
		return database.User{}, err
	}
	return user, nil
}

func (userService *UserService) EmailLogin(req request.Login) (database.User, error) {
	var user database.User
	err := global.DB.Where("email = ?", req.Email).First(&user).Error
	if err == nil {
		if ok := utils.BcryptCheck(req.Password, user.Password); !ok {
			return database.User{}, errors.New("incorrect email or password")
		}
		return user, nil
	}
	return database.User{}, err
}

func (userService *UserService) ForgotPassword(req request.ForgotPassword) error {
	user := database.User{}
	if err := global.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return err
	}
	user.Password = utils.BcryptHash(req.NewPassword)
	return global.DB.Save(&user).Error
}

func (userService *UserService) UserCard(req request.UserCard) (response.UserCard, error) {
	var user database.User
	if err := global.DB.Where("uuid = ?", req.UUID).Select("uuid", "username", "avatar", "address", "signature").First(&user).Error; err != nil {
		return response.UserCard{}, err
	}
	return response.UserCard{
		UUID:      user.UUID,
		Username:  user.Username,
		Avatar:    user.Avatar,
		Address:   user.Address,
		Signature: user.Signature,
	}, nil
}

func (userService *UserService) Logout() {}

func (userService *UserService) UserResetPassword() {}

func (userService *UserService) UserInfo() {}

func (userService *UserService) UserChangeInfo() {}

func (userService *UserService) UserWeather() {}

func (userService *UserService) UserChart() {}

func (userService *UserService) UserList() {}

func (userService *UserService) UserFreeze() {}

func (userService *UserService) UserUnfreeze() {}

func (userService *UserService) UserLoginList() {}
