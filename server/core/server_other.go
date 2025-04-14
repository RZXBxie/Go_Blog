package core

// 只有当系统不是Windows才用到这个方法
//func initServer(addr string, router *gin.Engine) server {
//	s := endless.NewServer(addr, router)
//	s.ReadHeaderTimeout = time.Minute * 10
//	s.WriteTimeout = time.Minute * 10
//	s.MaxHeaderBytes = 1 << 20
//
//	return s
//}
