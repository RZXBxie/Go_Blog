package service

type ServiceGroup struct {
	BaseService
	EsService
	JwtService
	GaodeService
	UserService
	QQService
	ImageService
}

var ServiceGroupApp = new(ServiceGroup)
