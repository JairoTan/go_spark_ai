package main

import (
	"fmt"
	"log"
	"net/http"

	"wxcloudrun-golang/db"
	"wxcloudrun-golang/service"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}
	router := gin.Default()

	router.GET("/message/post", service.SparkAIHandler)

	fmt.Println("Service Listening port: 80")
	log.Fatal(http.ListenAndServe(":80", nil))
}
