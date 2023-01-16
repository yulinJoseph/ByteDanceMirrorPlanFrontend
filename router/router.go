package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"im-instance/service"
	"im-instance/utils"
	"net/http"
	"strings"
)

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		origin := context.Request.Header.Get("Origin")
		var headerKeys []string
		for k := range context.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ",")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			context.Header("Access-Control-Allow-Origin", "*") // 设置允许访问所有域
			context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			context.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
			context.Header("Access-Control-Max-Age", "172800")
			context.Header("Access-Control-Allow-Credentials", "false")
			context.Set("content-type", "application/json") //// 设置返回格式是json
		}
		if method == "OPTIONS" {
			context.JSON(http.StatusOK, "Options Request!")
		}
		//处理请求
		context.Next()
	}
}

func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if utils.CheckToken(c.GetHeader("Authorization")) {
			c.Next()
		} else {
			c.JSON(http.StatusFound, gin.H{
				"message": "login first",
			})
			c.Abort()
		}
	}
}

func Router() *gin.Engine {
	r := gin.Default()
	r.Use(Cors())
	r.GET("/validateToken", service.ValidateToken)      // 验证 token
	r.POST("/register", service.Register)               // 注册用户
	r.POST("/login", service.Login)                     // 用户登录
	r.GET("/user/checkUsername", service.CheckUsername) // 检查用户名是否重复
	r.GET("/user/checkEmail", service.CheckEmail)       // 检查邮箱是否重复
	r.GET("/getAvatar", service.GetAvatar)              // 获取头像

	r.GET("/sendMsg", service.SendMsg) // websocket 建立连接

	r.Use(CheckToken())
	{
		r.GET("/user/getSelfInfo", service.GetSelfInfo)                       // 根据 token 获取自身信息
		r.GET("/user/getContactList", service.GetContactList)                 // 根据 token 获取自己的好友
		r.GET("/user/getMessages", service.GetMessages)                       // 根据 token 获取自己和对应好友的消息
		r.GET("/user/getContactWithMessages", service.GetContactWithMessages) // 根据 token 获取自己的好友和消息
		r.GET("/user/getUserList", service.GetUserList)                       // 获取到所有用户信息
		r.POST("/uploadAvatar", service.UploadAvatar)                         // 上传头像
		r.GET("/user/updateUsername", service.UpdateUsername)                 // 更新用户名
		r.GET("/user/updateSignature", service.UpdateSignature)               // 更新个性签名
		r.GET("/user/sendInvitation", service.SendInvitation)                 // 发送好友邀请
		r.GET("/user/getUserByID", service.GetUserByID)                       // 根据 id 获取用户信息
		r.GET("/agreeInvitation", service.AgreeInvitation)                    // 接受好友邀请
		r.GET("/disagreeInvitation", service.DisagreeInvitation)              // 接受好友邀请
	}

	return r
}
