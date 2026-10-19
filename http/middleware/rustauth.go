package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/global"
	"github.com/lejianwen/rustdesk-api/v2/service"
)

func RustAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//fmt.Println(c.Request.URL, c.Request.Header)
		//get HTTP_AUTHORIZATION
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}
		if len(token) <= 7 {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}
		//extract the token, the format is Bearer {token}
		//just a simple extraction here
		token = token[7:]

		//verify the token

		//check whether a jwt key is set
		if len(global.Jwt.Key) > 0 {
			uid, _ := service.AllService.UserService.VerifyJWT(token)
			if uid == 0 {
				c.JSON(401, gin.H{
					"error": "Unauthorized",
				})
				c.Abort()
				return
			}
		}

		user, ut := service.AllService.UserService.InfoByAccessToken(token)
		if user.Id == 0 {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}
		if !service.AllService.UserService.CheckUserEnable(user) {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		c.Set("curUser", user)
		c.Set("token", token)

		service.AllService.UserService.AutoRefreshAccessToken(ut)

		c.Next()
	}
}
