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

func useUploadFile(profile *model.Profile, c *gin.Context) error {
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
		if profile.Profile != "" {
			useRemoveFile(profile.Profile)
		}
		// Assign file value to model
		profile.Profile = fileName
		return nil
	}
	// Remove profile to empty
	reset := c.Query("reset")
	if profile.Profile != "" && reset == "1" {
		useRemoveFile(profile.Profile)
		// Assign empty value to model
		profile.Profile = ""
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

func UpdateProfile(c *gin.Context) {
	profile := model.Profile{}

	// Validate Pipe
	if err := validationPipeId(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	if err := validationPipeIdNotEqual(c.Param("id"), c.PostForm("user_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	if err := validationPipeIdNotEqual(c.Param("id"), c.PostForm("user_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	// Create Form
	findProfile := db.GetDB().Where("user_id = ?", c.Param("id")).First(&profile)
	// Err
	if findProfile.Error != nil {
		println(findProfile.Error)
	}
	if findProfile.RowsAffected == 0 {
		// Create new Profile
		userId, _ := strconv.ParseInt(c.PostForm("user_id"), 10, 32)
		profile.UserId = uint(userId)
		// Use Upload File
		if err := useUploadFile(&profile, c); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}
		// Create now!
		if err := db.GetDB().Save(&profile).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err})
			return
		}
		// Response
		c.JSON(http.StatusNotFound, gin.H{"success": true, "data": profile})
	} else {
		// Use Upload File
		if err := useUploadFile(&profile, c); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}
		// Update now!
		if err := db.GetDB().Model(&profile).Update("profile", profile.Profile).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err})
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": profile})
	}
}

func GetProfileById(c *gin.Context) {
	profile := model.Profile{}
	// Invalid id
	if err := validationPipeId(c.Param("id")); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	// Response
	if err := db.GetDB().Where("id = ?", c.Param("id")).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": profile})
}
