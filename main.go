//package main
//
//import (
//	"fmt"
//	"log"
//
//	"go_spark_ai/db"
//	"go_spark_ai/service"
//
//	"github.com/gin-gonic/gin"
//)
//
//func main() {
//	if err := db.Init(); err != nil {
//		panic(fmt.Sprintf("mysql init failed with %+v", err))
//	}
//	router := gin.Default()
//
//	router.POST("/message/post", service.SparkAIHandler)
//
//	fmt.Println("Service Listening port: 80")
//	log.Fatal(router.Run(":80"))
//}

package main

import (
	"net/http"

	"go_spark_ai/util"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/send", func(c *gin.Context) {
		var reqMsg struct {
			ToUserName   string `json:"ToUserName"`
			FromUserName string `json:"FromUserName"`
			MsgType      string `json:"MsgType"`
			Content      string `json:"Content"`
			CreateTime   int64  `json:"CreateTime"`
		}

		if err := c.BindJSON(&reqMsg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON数据包解析失败"})
			return
		}

		if reqMsg.MsgType == "text" {
			switch reqMsg.Content {
			case "你是":
				//获取星火AI
				answer := util.SparkAnswer(reqMsg.Content)
				c.JSON(http.StatusOK, gin.H{
					"ToUserName":   reqMsg.FromUserName,
					"FromUserName": reqMsg.ToUserName,
					"CreateTime":   reqMsg.CreateTime,
					"MsgType":      "text",
					"Content":      answer,
				})
			default:
				c.String(http.StatusOK, "success")
			}
		} else {
			c.String(http.StatusOK, "success")
		}
	})

	r.Run(":80")
}
