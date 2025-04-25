package service

type ServiceGroup struct {
	BaseService
	EsService
	JwtService
}

var ServiceGroupApp = new(ServiceGroup)
