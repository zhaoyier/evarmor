package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
)

//请求信息转换
func GinBindJson(c *gin.Context, req proto.Message) {
	// n := new(req)
	// c.BindJSON()

}
