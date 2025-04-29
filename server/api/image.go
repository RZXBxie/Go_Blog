package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"server/global"
	"server/model/request"
	"server/model/response"
)

type ImageApi struct{}

// ImageUpload 上传图片
func (imageApi *ImageApi) ImageUpload(ctx *gin.Context) {
	_, header, err := ctx.Request.FormFile("image")
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	url, err := imageService.ImageUpload(header)
	if err != nil {
		global.Log.Error("Failed to upload the image:", zap.Error(err))
		response.FailWithMessage("Failed to upload the image", ctx)
		return
	}
	response.OkWithDetailed(response.ImageUpload{
		Url:     url,
		OssType: global.Config.System.OssType,
	}, "Successfully upload the image", ctx)
}
func (imageApi *ImageApi) ImageDelete(ctx *gin.Context) {
	var req request.ImageDelete
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	err := imageService.ImageDelete(req)
	if err != nil {
		global.Log.Error("Failed to delete image:", zap.Error(err))
		response.FailWithMessage("Failed to delete image", ctx)
		return
	}
	response.OkWithMessage("Successfully delete the image", ctx)
}
func (imageApi *ImageApi) ImageList(ctx *gin.Context) {
	var req request.ImageList
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	list, total, err := imageService.ImageList(req)
	if err != nil {
		global.Log.Error("Failed to get image list:", zap.Error(err))
		response.FailWithMessage("Failed to get image list", ctx)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:  list,
		Total: total,
	}, "Successfully get image list", ctx)
}
