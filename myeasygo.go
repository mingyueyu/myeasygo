package myeasygo

// package main

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/mingyueyu/myeasygo/mmysql"
// 	"github.com/mingyueyu/myeasygo/util/mysqlTool"
// )

// type Options struct {
// 	Origin string
// }

// func main() {
// 	setupMysql()
// 	r := gin.Default()
// 	r.Use(CORS(Options{Origin: "*"}))
// 	r.Use(gin.Recovery())
// 	mmysql.Add(r, "/api/add", "test", "test", false, false, false)
// 	mmysql.Delete(r, "/api/delete", "test", "test")
// 	mmysql.Update(r, "/api/update", "test", "test")
// 	mmysql.List(r, "/api/list", "test", "test", []string{"name", "age"})
// 	mmysql.Detail(r, "/api/detail", "test", "test")
// 	mmysql.Dif(r, "/api/dif", "test", "test")
// 	r.Run(":12345")
// }


// func setupMysql(){
// 	mysqlTool.MysqlToolInit([]mysqlTool.MySql_t{
// 		{
// 			Host:   "localhost",
// 			Name:   "test",
// 			Port:   3306,
// 			Pwd:    "12345678",
// 			User:   "root",
// 			Tables: []mysqlTool.Table_t{
// 				{
// 					Name:    "test",
// 					Content: "name varchar(255), age varchar(255)",
// 				},
// 			},
// 		},
// 	})
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