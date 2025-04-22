package response

type CaptchaResponse struct {
	CaptchaID   string `json:"captcha_id"`
	PicturePath string `json:"picture_path"`
}
