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

func uuid() string {
	var chars = []rune("0123456789")
	str := make([]rune, 10)
	for i := range str {
		str[i] = chars[rand.Intn(len(chars))]
	}
	return string(str)
}

func useRemoveFile(file string) error {
	info, err := os.Stat(file)
	if !info.IsDir() && err == nil {
		err := os.Remove(file)
		if err != nil {
			return err
		}
	}
	return nil
}

type error interface {
	Error() string
}

func useUploadFile(form *model.Avatar, c *gin.Context) (error, string) {
	dir, _ := os.Getwd()
	file, _ := c.FormFile("file")
	if file != nil {
		// Check file
		fileSize := file.Size
		maxFileSize := 4194304 // 4MB
		if fileSize > int64(maxFileSize) {
			return errors.New("max file size 4mb"), ""
		}
		// Upload the file to specific dst.
		uniqueId := uuid()
		fileTmp := fmt.Sprintf("%s/tmp/%s%s", dir, uniqueId, ".tmp")
		saveUploadedErr := c.SaveUploadedFile(file, fileTmp)
		if saveUploadedErr != nil {
			log.Println(saveUploadedErr)
		}
		// Compress to Webp
		fileWebp := fmt.Sprintf("%s/public/%s%s", dir, uniqueId, ".webp")
		execWebp := exec.Command("cwebp", "-resize", "360", "360", "-o", fileWebp, fileTmp)
		execWebpErr := execWebp.Run()
		if execWebpErr != nil {
			return errors.New("err"), "error"
		}
		// Remove original file
		fileOrg := fmt.Sprintf("%s/public/%s", dir, form.Avatar)
		if form.Avatar != "" {
			if err := useRemoveFile(fileOrg); err != nil {
				log.Println(err)
				return err, ""
			}
		}

		// Remove tmp file
		// defer useRemoveFile(fileTmp)

		// if err, compress := compressWebp(fileId, fileExt, fileOrg); err {
		// 	return true, compress
		// }

		return nil, uniqueId + ".webp"
	}
	// Remove avatar to empty
	if form.Avatar != "" {
		fileOrg := fmt.Sprintf("%s/public/%s", dir, form.Avatar)
		if err := useRemoveFile(fileOrg); err != nil {
			log.Println(err)
			return err, ""
		}
		return nil, ""
	}
	return nil, ""
}

func validationPipeId(c *gin.Context) (bool, string) {
	if id, _ := strconv.ParseInt(c.Param("id"), 10, 32); id == 0 {
		return true, "invalid id"
	}
	return false, ""
}

func validationPipeIdNotEqual(c *gin.Context) (bool, string) {
	if c.Param("id") != c.PostForm("user_id") {
		return true, "id should be equal"
	}
	return false, ""
}

func UpdateAvatar(c *gin.Context) {
	// invalid id
	if errId, validateId := validationPipeId(c); errId {
		c.JSON(http.StatusBadRequest, gin.H{"message": validateId})
		return
	}
	if errIdNotEqual, validateIdNotEqual := validationPipeIdNotEqual(c); errIdNotEqual {
		c.JSON(http.StatusBadRequest, gin.H{"message": validateIdNotEqual})
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
		errFile, fileResult := useUploadFile(&avatar, c)
		if errFile != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": fileResult})
			return
		}
		avatar.Avatar = fileResult
		// Create now!
		if err := db.GetDB().Save(&avatar).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err})
			return
		}
		// Response
		c.JSON(http.StatusNotFound, gin.H{"success": true, "data": avatar})
	} else {
		// Use Upload File
		errFile, fileResult := useUploadFile(&avatar, c)
		if errFile != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": fileResult})
			return
		}
		avatar.Avatar = fileResult
		// Update now!
		if err := db.GetDB().Model(&avatar).Update("avatar", avatar.Avatar).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": avatar})
	}
}

func GetAvatarById(c *gin.Context) {
	avatar := model.Avatar{}
	// invalid id
	if errId, validateId := validationPipeId(c); errId {
		c.JSON(http.StatusBadRequest, gin.H{"message": validateId})
		return
	}
	// Response
	if err := db.GetDB().Where("id = ?", c.Param("id")).First(&avatar).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": avatar})
}
