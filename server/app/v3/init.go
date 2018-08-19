package v3

import (
	"github.com/gin-gonic/gin"
	"github.com/xuebing1110/rtbus/server/app"
)

var (
	router *gin.RouterGroup = app.Get().Group("/api/v3")
)
