package api

import (
	"errors"
	"fmt"
	"log"
	"main/db"
	"main/model"
	"path/filepath"
	"strings"

	"strconv"

	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type error interface {
	Error() string
}

func useRemoveFile(file string) bool {
	dir, _ := os.Getwd()
	if file == "" {
		return false
	}
	rm := fmt.Sprintf("%s/public/%s", dir, file)
	if err := os.Remove(rm); err != nil {
		log.Println(err)
	}
	return true
}

func useUploadFile(avatar *model.Avatar, c *gin.Context) error {
	dir, _ := os.Getwd()
	file, _ := c.FormFile("file")
	if file != nil {
		// Check file
		fileExt := filepath.Ext(file.Filename)
		fileSize := file.Size
		maxFileSize := 1048576 // 1MB
		if fileSize > int64(maxFileSize) {
			return errors.New("max file size 1mb")
		}
		// Upload the file to specific dst.
		uuidv4 := strings.Replace(uuid.New().String(), "-", "", -1)
		fileName := fmt.Sprintf("%v%s", uuidv4, fileExt)
		fileTmp := fmt.Sprintf("%s/public/%s", dir, fileName)
		if err := c.SaveUploadedFile(file, fileTmp); err != nil {
			log.Println(err)
		}
		// Remove original File
		if avatar.Avatar != "" {
			useRemoveFile(avatar.Avatar)
		}
		// Assign file name to model
		avatar.Avatar = fileName
		return nil
	}
	// Remove avatar to empty
	if avatar.Avatar != "" {
		fileOrg := fmt.Sprintf("%s/public/%s", dir, avatar.Avatar)
		useRemoveFile(fileOrg)
		return nil
	}
	return nil
}

func validationPipeId(id string) error {
	if id, _ := strconv.ParseInt(id, 10, 32); id == 0 {
		return errors.New("invalid id")
	}
	return nil
}

func validationPipeIdNotEqual(id string, user_id string) error {
	if id != user_id {
		return errors.New("id should be equal")
	}
	return nil
}

func UpdateAvatar(c *gin.Context) {
	// Invalid id
	if err := validationPipeId(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	if err := validationPipeIdNotEqual(c.Param("id"), c.PostForm("user_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	// Create Form
	avatar := model.Avatar{}
	findAvatar := db.GetDB().Where("user_id = ?", c.Param("id")).First(&avatar)
	// Err
	if findAvatar.Error != nil {
		println(findAvatar.Error)
	}

	if findAvatar.RowsAffected == 0 {
		// Create new Avatar
		userId, _ := strconv.ParseInt(c.PostForm("user_id"), 10, 32)
		avatar.UserId = uint(userId)
		// Use Upload File
		if err := useUploadFile(&avatar, c); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}
		// Create now!
		if err := db.GetDB().Save(&avatar).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err})
			return
		}
		// Response
		c.JSON(http.StatusNotFound, gin.H{"success": true, "data": avatar})
	} else {
		// Use Upload File
		if err := useUploadFile(&avatar, c); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}
		// Update now!
		if err := db.GetDB().Model(&avatar).Update("avatar", avatar.Avatar).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err})
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": avatar})
	}

}

func GetAvatarById(c *gin.Context) {
	avatar := model.Avatar{}
	// Invalid id
	if err := validationPipeId(c.Param("id")); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	// Response
	if err := db.GetDB().Where("id = ?", c.Param("id")).First(&avatar).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": avatar})
}
