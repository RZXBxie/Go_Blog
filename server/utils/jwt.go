package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"server/global"
	"server/model/request"
	"time"
)

type JWT struct {
	AccessTokenSecret  []byte
	RefreshTokenSecret []byte
}

var (
	TokenExpired     = errors.New("token is Expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")    // token格式错误
	TokenInvalid     = errors.New("couldn't handle this token") // token无效
)

// NewJWT 创建一个新的 JWT 实例，初始化 AccessToken 和 RefreshToken 密钥
func NewJWT() *JWT {
	return &JWT{
		AccessTokenSecret:  []byte(global.Config.Jwt.AccessTokenSecret),
		RefreshTokenSecret: []byte(global.Config.Jwt.RefreshTokenSecret),
	}
}

// CreateAccessClaims 创建 Access Token 的 Claims，包含基本信息和过期时间等
func (j *JWT) CreateAccessClaims(baseClaims request.BaseClaims) request.JwtCustomClaims {
	expiration, _ := ParseDuration(global.Config.Jwt.AccessTokenExpiryTime)
	claims := request.JwtCustomClaims{
		BaseClaims: baseClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"Go Blog"},                    // 受众
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)), // 过期时间
			Issuer:    global.Config.Jwt.Issuer,                       // 签名的发起者
		},
	}
	return claims
}

func (j *JWT) CreateAccessToken(claims request.JwtCustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.AccessTokenSecret) // 使用accessTokenSecret进行签名
}

func (j *JWT) CreateRefreshClaims(baseClaims request.BaseClaims) request.JwtCustomRefreshClaims {
	expiration, _ := ParseDuration(global.Config.Jwt.RefreshTokenExpiryTime)
	claims := request.JwtCustomRefreshClaims{
		UserID: baseClaims.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"Go Blog"},                    // 受众
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)), // 过期时间
			Issuer:    global.Config.Jwt.Issuer,
		},
	}
	return claims
}

func (j *JWT) CreateRefreshToken(claims request.JwtCustomRefreshClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.RefreshTokenSecret)
}

// ParseAccessToken 解析 Access Token，验证 Token 并返回 Claims 信息
func (j *JWT) ParseAccessToken(tokenString string) (*request.JwtCustomClaims, error) {
	claims, err := j.parseToken(tokenString, &request.JwtCustomClaims{}, global.Config.Jwt.AccessTokenSecret)
	if err != nil {
		return nil, err
	}
	if customClaims, ok := claims.(*request.JwtCustomClaims); ok { // 确保解析出的claims正确
		return customClaims, nil
	}
	return nil, TokenInvalid // 如果解析无果，返回TokenInvalid错误
}

// ParseRefreshToken 解析 Refresh Token，验证 Token 并返回 Claims 信息
func (j *JWT) ParseRefreshToken(tokenString string) (*request.JwtCustomRefreshClaims, error) {
	claims, err := j.parseToken(tokenString, &request.JwtCustomRefreshClaims{}, j.RefreshTokenSecret) // 解析 Token
	if err != nil {
		return nil, err
	}
	if refreshClaims, ok := claims.(*request.JwtCustomRefreshClaims); ok { // 确保解析出的 Claims 类型正确
		return refreshClaims, nil
	}
	return nil, TokenInvalid // 如果解析结果无效，返回 TokenInvalid 错误
}

// parseToken 通用的 Token 解析方法，验证 Token 是否有效并返回 Claims
func (j *JWT) parseToken(tokenString string, claims jwt.Claims, secretKey interface{}) (interface{}, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			switch {
			case ve.Errors&jwt.ValidationErrorMalformed != 0:
				return nil, TokenMalformed // Token 格式错误
			case ve.Errors&jwt.ValidationErrorExpired != 0:
				return nil, TokenExpired // Token 已过期
			case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
				return nil, TokenNotValidYet // Token 还未生效
			default:
				return nil, TokenInvalid // 其他错误返回 Token 无效
			}
		}
		return nil, TokenInvalid
	}
	if token.Valid {
		return token.Claims, nil
	}
	return nil, TokenInvalid
}
