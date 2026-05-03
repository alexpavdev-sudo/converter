package controllers

import (
	"converter/app"
	"converter/config"
	"converter/dto/web"
	"converter/entities"
	"converter/helpers"
	deleteFile "converter/services/delete"
	"converter/services/uploader"
	"converter/services/user"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

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
		guestId, err := userService.InitGuestID()
		if err != nil {
			app.Fail(c, 400, "invalid_guest", "Guest not found")
			return
		}

		files, err := app.App().FileRepo.GetFiles(guestId)
		if err != nil {
			app.Fail(c, 500, "1", err.Error())
			return
		}
		app.OK(c, web.FileToDTO().MapSlice(files))
	}
}

func GetFile(c *gin.Context) {
	file, err := findFile(c)
	if err != nil {
		app.Fail(c, 500, "1", err.Error())
		return
	}

	app.OK(c, web.FileToDTO().Map(*file))
}

func findFile(c *gin.Context) (*entities.File, error) {
	var req web.FileRequest
	if err := c.ShouldBindUri(&req); err != nil {
		return nil, err
	}

	userService := user.NewUserService(sessions.Default(c))
	guestId, err := userService.InitGuestID()
	if err != nil {
		return nil, errors.New("guest not found")
	}
	file, err := app.App().FileRepo.GetFile(guestId, req.ID)
	if err != nil || file.ID <= 0 {
		return nil, errors.New("file not found")
	}

	return &file, nil
}

func GetFileError(c *gin.Context) {
	file, err := findFile(c)
	if err != nil {
		app.Fail(c, 500, "1", err.Error())
		return
	}

	var errorModel entities.Error
	err = app.App().DB.Model(&entities.Error{}).
		Where("file_id = ?", file.ID).
		Order("created_at DESC").
		Limit(1).
		First(&errorModel).Error

	if err != nil {
		app.Fail(c, 500, "1", err.Error())
		return
	}

	app.OK(c, errorModel.Details)
}

func DownloadFile(c *gin.Context) {
	file, err := findFile(c)
	if err != nil {
		app.Fail(c, 500, "1", err.Error())
		return
	}

	if file.Status != entities.StatusProcessed {
		app.Fail(c, 404, "1", "FileDto not processed")
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", helpers.GetFileNameWithoutExt(file.OriginalName)+"."+file.Format))
	c.File(file.ProcessedPathFull())
}

func DeleteFile(c *gin.Context) {
	file, err := findFile(c)
	if err != nil {
		app.Fail(c, 500, "1", err.Error())
		return
	}

	err = deleteFile.DeleteFile(*file)
	if err != nil {
		app.Fail(c, 500, "1", "Failed to delete file")
		return
	}

	app.OK(c, nil)
}
