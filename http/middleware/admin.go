package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/http/response"
	"github.com/lejianwen/rustdesk-api/v2/service"
)

// BackendUserAuth admin-backend permission middleware
func BackendUserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		//disabled for testing first
		token := c.GetHeader("api-token")
		if token == "" {
			response.Fail(c, 403, response.TranslateMsg(c, "NeedLogin"))
			c.Abort()
			return
		}
		user, ut := service.AllService.UserService.InfoByAccessToken(token)
		if user.Id == 0 {
			response.Fail(c, 403, response.TranslateMsg(c, "NeedLogin"))
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
		//if less than 1 day is left, renew the token automatically
		service.AllService.UserService.AutoRefreshAccessToken(ut)

		c.Next()
	}
}
