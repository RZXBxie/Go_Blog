package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"server/global"
	"server/model/apptypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/model/response"
	"server/utils"
	"time"
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

// Logout 用户登出
func (userService *UserService) Logout(c *gin.Context) {
	jwtStr := utils.GetRefreshToken(c)
	UUID := utils.GetUUID(c)
	global.Redis.Del(UUID.String())
	utils.ClearRefreshToken(c)
	_ = ServiceGroupApp.JwtService.InsertIntoBlacklist(database.JWTBlacklist{Jwt: jwtStr})
	
}

func (userService *UserService) UserResetPassword(req request.UserResetPassword) error {
	var user database.User
	if err := global.DB.Take(&user, req.UserID).Error; err != nil {
		return err
	}
	if ok := utils.BcryptCheck(req.Password, user.Password); !ok {
		return errors.New("incorrect original password")
	}
	user.Password = utils.BcryptHash(req.NewPassword)
	return global.DB.Save(&user).Error
}

func (userService *UserService) UserInfo(userID uint) (database.User, error) {
	var user database.User
	if err := global.DB.Take(&user, userID).Error; err != nil {
		return database.User{}, err
	}
	return user, nil
}

func (userService *UserService) UserChangeInfo(req request.UserChangeInfo) error {
	var user database.User
	if err := global.DB.Take(&user, req.UserID).Error; err != nil {
		return err
	}
	return global.DB.Model(&user).Updates(req).Error
}

func (userService *UserService) UserWeather(ip string) (string, error) {
	// 从redis中获取天气数据，如果没有，则利用高德api进行查询
	result, err := global.Redis.Get("weather-" + ip).Result()
	if err != nil {
		ipResponse, err := ServiceGroupApp.GaodeService.GetLocationByIP(ip)
		if err != nil {
			return "", err
		}
		live, err := ServiceGroupApp.GaodeService.GetWeatherByAdcode(ipResponse.Adcode)
		if err != nil {
			return "", err
		}
		weather := "地区：" + live.Province + "-" + live.City + " 天气：" + live.Weather + " 温度：" + live.Temperature + "°C" + " 风向：" + live.WindDirection + " 风级：" + live.WindPower + " 湿度：" + live.Humidity + "%"
		// 将天气数据存入redis
		if err := global.Redis.Set("weather-"+ip, weather, time.Hour*1).Err(); err != nil {
			return "", err
		}
		return weather, nil
	}
	return result, nil
}

// UserChart TODO 存疑，这里获取的是所有用户的登录和注册数据，不符合实际场景
func (userService *UserService) UserChart(req request.UserChart) (response.UserChart, error) {
	var resp response.UserChart
	// 构建查询条件
	where := global.DB.Where("date_sub(curdate(), interval ? day) <= created_at", req.Date)
	
	// 生成日期列表
	startDate := time.Now().AddDate(0, 0, -req.Date)
	for i := 1; i <= req.Date; i++ {
		resp.DateList = append(resp.DateList, startDate.AddDate(0, 0, i).Format("2006-01-02"))
	}
	// 获取登录数据
	loginData := utils.FetchDateCounts(global.DB.Model(&database.Login{}), where)
	
	// 获取注册数据
	registerData := utils.FetchDateCounts(global.DB.Model(&database.User{}), where)
	
	for _, date := range resp.DateList {
		loginCount := loginData[date]
		registerCount := registerData[date]
		resp.LoginData = append(resp.LoginData, loginCount)
		resp.RegisterData = append(resp.RegisterData, registerCount)
	}
	return resp, nil
}

func (userService *UserService) UserList(req request.UserList) (interface{}, int64, error) {
	db := global.DB
	if req.UUID != nil {
		db = db.Where("uuid = ?", req.UUID)
	}
	if req.Register != nil {
		db = db.Where("register = ?", req.Register)
	}
	
	option := other.MySQLOption{
		PageInfo: req.PageInfo,
		Where:    db,
	}
	return utils.MySQLPagination(&database.User{}, option)
}

func (userService *UserService) UserFreeze(req request.UserFreeze) error {
	var user database.User
	jwtService := ServiceGroupApp.JwtService
	if err := global.DB.Take(&user, req.UserID).Update("freeze", true).Error; err != nil {
		return err
	}
	jwtStr, _ := jwtService.GetJwtFromRedis(user.UUID)
	if jwtStr != "" {
		_ = jwtService.InsertIntoBlacklist(database.JWTBlacklist{Jwt: jwtStr})
	}
	return nil
}

func (userService *UserService) UserUnfreeze(req request.UserFreeze) error {
	return global.DB.Take(&database.User{}, req.UserID).Update("freeze", false).Error
}

func (userService *UserService) UserLoginList(req request.UserLoginList) (interface{}, int64, error) {
	db := global.DB
	if req.UUID != nil {
		var userID uint
		if err := global.DB.Model(database.User{}).Where("uuid = ?", req.UUID).Pluck("id", &userID).Error; err != nil {
			return nil, 0, err
		}
		db.Where("user_id = ?", userID)
	}
	option := other.MySQLOption{
		PageInfo: req.PageInfo,
		Where:    db,
		Preload:  []string{"User"},
	}
	return utils.MySQLPagination(&database.Login{}, option)
}
