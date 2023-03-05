package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"im-instance/models"
	"im-instance/utils"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

func GetUserList(c *gin.Context) {
	if data, err := models.GetUserList(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
	} else {
		if data != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": data,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "no user exists",
			})
		}
	}
}

func CheckUsername(c *gin.Context) {
	fmt.Println(c.Query("username"))
	if models.CheckUsername(c.Query("username")) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "username exists",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "username not exists",
		})
	}
}

func CheckEmail(c *gin.Context) {
	if models.CheckEmail(c.Query("email")) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "email exists",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "email not exists",
		})
	}
}

func Register(c *gin.Context) {
	user := models.UserBasic{
		Model: gorm.Model{
			ID:        0,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
	}
	if err := c.ShouldBind(&user); err != nil {
		fmt.Println("asd")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("%#v\n", user)
	if err := models.CreateUser(&user); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "create error",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "created successfully",
		})
	}
}

func Login(c *gin.Context) {
	var user models.UserBasic
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if models.CheckUserByUsername(&user) || models.CheckUserByEmail(&user) {
		token, err := utils.NewToken(fmt.Sprintf("%d+%s+%s", user.ID, user.Username, user.Passwd))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "user exists",
			"token":   token,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("username or passwd not correct"),
		})
	}
}

func GetSelfInfo(c *gin.Context) {
	id := utils.GetIdFromToken(c.GetHeader("Authorization"))
	user, err := models.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": user,
	})
}

func SendMsg(c *gin.Context) {
	token := c.Query("token")
	fromID := utils.GetIdFromToken(token)
	if fromID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token invalid",
		})
	}

	models.Chat(c.Writer, c.Request, uint64(fromID))
}

func UploadAvatar(c *gin.Context) {
	id := utils.GetIdFromToken(c.GetHeader("Authorization"))
	if id == 0 {
		fmt.Println("token invalid")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token invalid",
		})
	}
	file, err := c.FormFile("avatar")
	if err != nil {
		fmt.Println("get file error")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "upload image error",
		})
		return
	}
	ext := strings.ToLower(path.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		fmt.Println("file type error")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "image format not supported",
		})
		return
	}

	if err := c.SaveUploadedFile(file, fmt.Sprintf("public/avatar/%d%s", id, ext)); err != nil {
		fmt.Println("save file error")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "save image error",
		})
		return
	}
	if err := models.UploadAvatar(id, fmt.Sprintf("%d%s", id, ext)); err != nil {
		fmt.Println("UploadAvatar() error")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "update avatar error",
		})
		return
	}
	defer c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "upload avatar successfully",
	})
}

func GetAvatar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "id error",
		})
		return
	}
	avatar, err := models.GetAvatar((uint)(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "get avatar error",
		})
		return
	}

	c.File(fmt.Sprintf("./public/avatar/%s", avatar))
}

func UpdateUsername(c *gin.Context) {
	id := utils.GetIdFromToken(c.GetHeader("Authorization"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token invalid",
		})
	}
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "username is empty",
		})
	}
	if err := models.UpdateUsername(id, username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "update username error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "update username successfully",
	})
}

func UpdateSignature(c *gin.Context) {
	id := utils.GetIdFromToken(c.GetHeader("Authorization"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token invalid",
		})
	}
	signature := c.Query("signature")
	if err := models.UpdateSignature(id, signature); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "update signature error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "update signature successfully",
	})
}

func GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "id error",
		})
	}
	user, err := models.GetUserByID((uint)(id))
	user.Passwd = ""
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "user not exists",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": user,
	})
}
