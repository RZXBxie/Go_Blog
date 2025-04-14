package core

import (
	"go.uber.org/zap"
	"server/global"
	"server/initialize"
)

type server interface {
	ListenAndServe() error
}

func RunServer() {
	addr := global.Config.System.Addr()
	Router := initialize.InitRouter()

	// TODO 加载所有的JWT黑名单，并存入本地缓存

	// 初始化服务器并启动
	s := initServer(addr, Router)
	global.Log.Info("server run success on ", zap.String("addr:", addr))
	global.Log.Error(s.ListenAndServe().Error())
}
