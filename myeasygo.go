package myeasygo

// package main

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/mingyueyu/myeasygo/demo"
// 	"github.com/mingyueyu/myeasygo/util/system"
// )

// type Options struct {
// 	Origin string
// }

// func main() {
// 	r := gin.Default()
// 	r.Use(CORS(Options{Origin: "*"}))
// 	r.Use(gin.Recovery())
// 	// 读取配置
// 	demo.SettingInit()
// 	// 数据库demo
// 	demo.MysqlInit(r)
// 	// tea算法
// 	system.TeaDemo()
// 	r.Run(":12345")
// }



// func CORS(options Options) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if options.Origin != "" {
// 			c.Writer.Header().Set("Access-Control-Allow-Origin", options.Origin)
// 		}
// 		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Connection, User-Agent, Cookie")
// 		c.Writer.Header().Set("Access-Control-Expose-Headers", "*")
// 		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(http.StatusNoContent)
// 		} else {
// 			c.Next()
// 		}
// 	}
// }