package main

import (
	"fmt"
	"log"

	"go_spark_ai/db"
	"go_spark_ai/service"

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
