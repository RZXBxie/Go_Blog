# 项目开发过程中遇到的问题合集
## 1.session存储数据出错
### 问题描述
使用session存储邮箱验证码过期时间时，错误的存入了time.Time类型，并且没有使用gob.Register进行注册，导致session存入失败
### 对应代码
```go
    package service
    expirationTime := time.Now().Add(5 * time.Minute)

	// 将验证码、发送邮箱和过期时间存入会话中
	session := sessions.Default(c)
	session.Set("verification_code", verificationCode)
	session.Set("expiration_time", expirationTime)
	session.Set("email", to)
	if err := session.Save(); err != nil {
		global.Log.Error("Failed to save session:", zap.Error(err))
	}
```
### 解决措施
存储简单类型（如int64），不存入time.Time类型
```go
expirationTime := time.Now().Add(5 * time.Minute).Unix()
```
