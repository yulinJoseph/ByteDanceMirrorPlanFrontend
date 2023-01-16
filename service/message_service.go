package service

import (
	"github.com/gin-gonic/gin"
	"im-instance/models"
	"im-instance/utils"
	"net/http"
	"strconv"
)

func GetMessages(c *gin.Context) {
	token := c.GetHeader("Authorization")
	fromID := utils.GetIdFromToken(token)
	toID, err := strconv.ParseUint(c.Query("toID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "toID invalid",
		})
		return
	}
	list, err := models.GetMsgList(fromID, uint(toID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "get messages failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": list,
	})
}
