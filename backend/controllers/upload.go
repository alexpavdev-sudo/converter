package controllers

import (
	"converter/app"
	"converter/config"
	"converter/dto"
	"converter/entities"
	"converter/helpers"
	deleteFile "converter/services/delete"
	"converter/services/uploader"
	"converter/services/user"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type FileRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

func Upload(c *gin.Context) {
	reader, err := c.Request.MultipartReader()
	if err != nil {
		app.Fail(c, 500, "1", "Failed to fetch files")
		return
	}

	session := sessions.Default(c)
	uploader := uploader.NewStreamFileUploader(reader, config.MaxSizeFile, config.MaxSize, session)
	err = uploader.Upload()
	if err != nil {
		app.Fail(c, 400, "1", err.Error())
		return
	}

	app.OK(c, gin.H{
		"count": uploader.CountSavedFiles(),
	})
}

func GetFiles(c *gin.Context) {
	userService := user.NewUserService(sessions.Default(c))

	if userService.IsAuthenticated() {
		_, err := userService.UserId()
		if err != nil {
			app.Fail(c, 500, "1", "Failed to get user id")
			return
		}
	} else {
		guestId, err := userService.GuestId()
		if err != nil {
			app.Fail(c, 400, "invalid_guest", "Guest not found")
			return
		}

		files, err := app.App().FileRepo.GetFiles(guestId)
		if err != nil {
			app.Fail(c, 500, "1", err.Error())
			return
		}
		app.OK(c, dto.FileToDTO().MapSlice(files))
	}
}

func GetFile(c *gin.Context) {
	var req FileRequest
	if err := c.ShouldBindUri(&req); err != nil {
		app.Fail(c, 400, "1", "Invalid file ID: "+err.Error())
		return
	}
	userService := user.NewUserService(sessions.Default(c))

	guestId, err := userService.GuestId()
	if err != nil {
		app.Fail(c, 400, "invalid_guest", "Guest not found")
		return
	}
	file, err := app.App().FileRepo.GetFile(guestId, req.ID)
	if err != nil {
		app.Fail(c, 500, "1", err.Error())
		return
	}
	app.OK(c, dto.FileToDTO().Map(file))
}

func DownloadFile(c *gin.Context) {
	var req FileRequest
	if err := c.ShouldBindUri(&req); err != nil {
		app.Fail(c, 400, "1", "Invalid file ID: "+err.Error())
		return
	}
	userService := user.NewUserService(sessions.Default(c))
	guestId, err := userService.GuestId()
	if err != nil {
		app.Fail(c, 400, "1", "failed to get guest id")
		return
	}
	file, err := app.App().FileRepo.GetFile(guestId, req.ID)
	if err != nil {
		app.Fail(c, 404, "1", "File not found")
		return
	}
	if file.Status != entities.StatusProcessed {
		app.Fail(c, 404, "1", "File not processed")
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", helpers.GetFileNameWithoutExt(file.OriginalName)+"."+file.Format))
	c.File(file.ProcessedPathFull())
}

func DeleteFile(c *gin.Context) {
	var req FileRequest
	if err := c.ShouldBindUri(&req); err != nil {
		app.Fail(c, 400, "1", "Invalid file ID: "+err.Error())
		return
	}
	userService := user.NewUserService(sessions.Default(c))
	guestId, err := userService.GuestId()
	if err != nil {
		app.Fail(c, 400, "1", "failed to get guest id")
		return
	}

	deleteService := deleteFile.NewDeleteService(sessions.Default(c))
	err = deleteService.DeleteFile(guestId, req.ID)
	if err != nil {
		app.Fail(c, 500, "1", "Failed to delete file")
		return
	}

	app.OK(c, nil)
}
