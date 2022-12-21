package session

import (
	"net/http"

	"github.com/lixvyang/betxin/internal/utils"

	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionToken := session.Get("token")
		if sessionToken == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Not logged",
			})
			c.Abort()
		}
		c.Next()
	}
}

func EnableCookileSession(r *gin.Engine) {
	store, _ := redisStore.NewStore(10, "tcp", utils.RedisHost+":"+utils.RedisPort, utils.RedisPassword, []byte(utils.RedisSecret))
	r.Use(sessions.Sessions("_betxin_session", store))
}
