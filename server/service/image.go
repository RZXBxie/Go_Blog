package service

import (
	"gorm.io/gorm"
	"mime/multipart"
	"server/global"
	"server/model/apptypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/upload"
	"server/utils"
)

type ImageService struct{}

func (imageService *ImageService) ImageUpload(header *multipart.FileHeader) (string, error) {
	oss := upload.NewOSS()
	url, filename, err := oss.UploadImage(header)
	if err != nil {
		return "", err
	}
	
	return url, global.DB.Create(&database.Image{
		Name:     filename,
		Path:     url,
		Storage:  global.Config.System.Storage(),
		Category: apptypes.Null,
	}).Error
}

func (imageService *ImageService) ImageDelete(req request.ImageDelete) error {
	if len(req.IDs) == 0 {
		return nil
	}
	var images []database.Image
	if err := global.DB.Find(&images, req.IDs).Error; err != nil {
		return err
	}
	for _, image := range images {
		if err := global.DB.Transaction(func(tx *gorm.DB) error {
			oss := upload.NewOssWithStorage(image.Storage)
			if err := global.DB.Delete(&image).Error; err != nil {
				return err
			}
			return oss.DeleteImage(image.Name)
		}); err != nil {
			return err
		}
	}
	return nil
}

func (imageService *ImageService) ImageList(req request.ImageList) (interface{}, int64, error) {
	db := global.DB
	
	if req.Name != nil {
		db = db.Where("name LIKE ?", "%"+*req.Name+"%")
	}
	
	if req.Category != nil {
		category := apptypes.ToCategory(*req.Category)
		db = db.Where("category = ?", category)
	}
	
	if req.Storage != nil {
		storage := apptypes.ToStorage(*req.Storage)
		db = db.Where("storage = ?", storage)
	}
	
	option := other.MySQLOption{
		Where:    db,
		PageInfo: req.PageInfo,
	}
	
	return utils.MySQLPagination(&database.Image{}, option)
}
