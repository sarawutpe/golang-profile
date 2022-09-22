package api

import (
	"errors"
	"fmt"
	"log"
	"main/db"
	"main/model"
	"math/rand"
	"os/exec"

	"strconv"

	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type error interface {
	Error() string
}

func uuid() string {
	var chars = []rune("0123456789")
	str := make([]rune, 10)
	for i := range str {
		str[i] = chars[rand.Intn(len(chars))]
	}
	return string(str)
}

func useRemoveFile(file string) {
	info, err := os.Stat(file)
	if !info.IsDir() && err == nil {
		err := os.Remove(file)
		if err != nil {
			log.Println(err)
		}
	}
}

func useUploadFile(avatar *model.Avatar, c *gin.Context) error {
	dir, _ := os.Getwd()
	file, _ := c.FormFile("file")
	if file != nil {
		// Check file
		fileSize := file.Size
		maxFileSize := 1048576 // 1MB
		if fileSize > int64(maxFileSize) {
			return errors.New("max file size 1mb")
		}
		// Remove Original File
		orgFileName := avatar.Avatar
		path := fmt.Sprintf("%s/public/%s", dir, orgFileName)
		if orgFileName != "" {
			exe := exec.Command("rm", ".", path)
			exe.Run()
		}

		// Upload the file to specific dst.
		uniqueId := uuid()
		fileTmp := fmt.Sprintf("%s/public/%s%s", dir, uniqueId, ".tmp")
		saveUploadedErr := c.SaveUploadedFile(file, fileTmp)
		if saveUploadedErr != nil {
			log.Println(saveUploadedErr)
		}
		// Assign file name to model
		avatar.Avatar = uniqueId + ".webp"
		return nil
	}
	// Remove avatar to empty
	if avatar.Avatar != "" {
		fileOrg := fmt.Sprintf("%s/public/%s", dir, avatar.Avatar)
		defer useRemoveFile(fileOrg)
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
	// invalid id
	if err := validationPipeId(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := validationPipeIdNotEqual(c.Param("id"), c.PostForm("user_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// Create Form
	avatar := model.Avatar{}
	findAvatar := db.GetDB().Where("user_id = ?", c.Param("id")).First(&avatar)
	if errors.Is(findAvatar.Error, gorm.ErrRecordNotFound) {
		// Create new Avatar
		userId, _ := strconv.ParseInt(c.PostForm("user_id"), 10, 32)
		avatar.UserId = uint(userId)
		// Use Upload File
		if err := useUploadFile(&avatar, c); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
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
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
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
	// invalid id
	if err := validationPipeId(c.Param("id")); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// Response
	if err := db.GetDB().Where("id = ?", c.Param("id")).First(&avatar).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": avatar})
}
