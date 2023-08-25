package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haobinfei/ginner/public/tools"
)

// 返回前端
func Response(c *gin.Context, httpStatus, code int, data gin.H, message string) {
	c.JSON(httpStatus, tools.H{
		"code":    code,
		"data":    data,
		"massage": message,
	})
}

func Success(c *gin.Context, data gin.H, message string) {
	Response(c, http.StatusOK, 200, data, message)
}
