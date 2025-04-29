package upload

import (
	"mime/multipart"
	"server/global"
	"server/model/apptypes"
)

var WhiteImageList = map[string]struct{}{
	".jpg":  {},
	".png":  {},
	".jpeg": {},
	".ico":  {},
	".tiff": {},
	".gif":  {},
	".svg":  {},
	".webp": {},
}

// OSS 对象存储接口定义，规定了文件上传和删除方法
type OSS interface {
	UploadImage(file *multipart.FileHeader) (string, string, error)
	DeleteImage(key string) error
}

// NewOSS 是OSS的实例化方法，根据配置中的OssType来决定存储类型
func NewOSS() OSS {
	switch global.Config.System.OssType {
	case "local":
		return &Local{}
	case "qiniu":
		return &Qiniu{}
	default:
		return &Local{}
	}
}

// NewOssWithStorage 是根据传入的存储类型返回相应的 OSS 实例
func NewOssWithStorage(storage apptypes.Storage) OSS {
	switch storage {
	case apptypes.Local:
		return &Local{}
	case apptypes.Qiniu:
		return &Qiniu{}
	default:
		return &Local{}
	}
}
