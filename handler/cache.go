package handler

import (
	"github.com/gin-gonic/gin"
)

// SetNoCacheHeaders 设置HTTP响应头，禁止缓存
func SetNoCacheHeaders(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
}

// HTMLWithNoCache 包装c.HTML方法，自动添加不缓存策略
func HTMLWithNoCache(c *gin.Context, code int, name string, obj interface{}) {
	SetNoCacheHeaders(c)
	c.HTML(code, name, obj)
}