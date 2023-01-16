package service

import (
	"github.com/gin-gonic/gin"
	"im-instance/utils"
)

func ValidateToken(c *gin.Context) {
	c.JSON(200, gin.H{
		"tokenValid": utils.CheckToken(c.GetHeader("Authorization")),
	})
}
