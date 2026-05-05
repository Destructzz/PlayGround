package pkg

import "github.com/gin-gonic/gin"

func SetGinMode(env string) {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
		return
	}

	gin.SetMode(gin.DebugMode)
}
