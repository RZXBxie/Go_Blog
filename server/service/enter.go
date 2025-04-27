package service

type ServiceGroup struct {
	BaseService
	EsService
	JwtService
	GaodeService
	UserService
	QQService
}

var ServiceGroupApp = new(ServiceGroup)
