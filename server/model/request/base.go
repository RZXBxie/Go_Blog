package request

type SendEmailVerificationCode struct {
	Email     string `json:"email" binding:"required,email"` // 前端必须传入这个参数，且会进行邮箱格式检验
	Captcha   string `json:"captcha" binding:"required,len=6"`
	CaptchaID string `json:"captcha_id" binding:"required"`
}
