package upload

import (
	"mime/multipart"
)

type Qiniu struct{}

func (q *Qiniu) UploadImage(file *multipart.FileHeader) (string, string, error) {
	//TODO implement me
	panic("implement me")
}

func (q *Qiniu) DeleteImage(key string) error {
	//TODO implement me
	panic("implement me")
}
