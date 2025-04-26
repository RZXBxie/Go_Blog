package core

import (
	"go.uber.org/zap"
	"server/global"
	"server/initialize"
	"server/service"
)

type server interface {
	ListenAndServe() error
}

func RunServer() {
	addr := global.Config.System.Addr()
	Router := initialize.InitRouter()

	// 加载所有JWT，并存入本地缓存
	service.LoadAll()

	// 初始化服务器并启动
	s := initServer(addr, Router)
	global.Log.Info("server run success on ", zap.String("addr:", addr))
	global.Log.Error(s.ListenAndServe().Error())
}
