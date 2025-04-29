package upload

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"server/global"
	"server/utils"
	"strings"
	"time"
)

type Local struct{}

func (l *Local) UploadImage(file *multipart.FileHeader) (string, string, error) {
	// 计算文件大小
	size := float64(file.Size) / float64(1024*1024)
	// 超过系统配置则报错
	if size >= float64(global.Config.Upload.Size) {
		return "", "", fmt.Errorf("the image size exceeds the set size, the current size is: %.2f MB, the set size is: %d MB", size, global.Config.Upload.Size)

	}
	// 获取文件扩展名（后缀）
	ext := filepath.Ext(file.Filename)
	name := strings.TrimSuffix(file.Filename, ext)
	if _, exists := WhiteImageList[ext]; !exists {
		return "", "", errors.New("don't upload files that aren't image types")
	}

	// 为了防止文件名冲突，先使用hash生成摘要
	filename := utils.MD5V([]byte(name)) + "-" + time.Now().Format("20060102150405") + ext
	path := global.Config.Upload.Path + "/image/"

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", "", err
	}

	// 文件路径=文件存储路径+文件名
	// uploads/image/******.png
	filePath := path + filename

	out, err := os.Create(filePath)
	if err != nil {
		return "", "", err
	}
	defer out.Close()

	f, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	if _, err = io.Copy(out, f); err != nil {
		return "", "", err
	}

	return "/" + filePath, filename, nil
}

func (l *Local) DeleteImage(key string) error {
	path := global.Config.Upload.Path + "/image/" + key
	return os.Remove(path)
}
