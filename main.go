package main

import (
	"fmt"
	"log"

	"wxcloudrun-golang/db"
	"wxcloudrun-golang/service"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}
	router := gin.Default()

	router.POST("/message/post", service.SparkAIHandler)

	fmt.Println("Service Listening port: 80")
	log.Fatal(router.Run(":80"))
}
