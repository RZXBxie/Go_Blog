package service

type ServiceGroup struct {
	BaseService
	EsService
}

var ServiceGroupApp = new(ServiceGroup)
