package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"go_spark_ai/util"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	//客服消息
	r.POST("/customer/message", func(c *gin.Context) {
		// 获取请求参数
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

		appid := c.GetHeader("x-wx-from-appid")
		openid := c.GetHeader("x-wx-openid")
		fmt.Printf("推送接收的账号 %s", reqMsg.FromUserName)
		fmt.Printf("公众号 %s 接收用户openid为 %s 的 %s 消息：%s", appid, openid, reqMsg.MsgType, reqMsg.Content)
		//获取星火AI
		var answer string
		if reqMsg.MsgType == "text" {
			answer, _ = util.SparkAnswer(reqMsg.Content)
			fmt.Println("星火AI回答：", answer)
		} else {
			answer = "暂不支持文本以外的消息回复"
			fmt.Println("不是文本消息,拒绝回答!")
		}

		// 构建消息体
		message := struct {
			ToUser  string `json:"touser"`
			MsgType string `json:"msgtype"`
			Text    struct {
				Content string `json:"content"`
			} `json:"text"`
		}{
			ToUser:  reqMsg.FromUserName,
			MsgType: "text",
		}
		message.Text.Content = reqMsg.Content

		// 将消息体转换为 JSON
		requestBody, err := json.Marshal(message)
		fmt.Println("requestBody为：", requestBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 发送 POST 请求到微信接口
		url := "http://api.weixin.qq.com/cgi-bin/message/custom/send?from_appid=" + appid
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
		fmt.Println("微信接口返回内容：", resp)
		if err != nil {
			fmt.Println("请求微信接口报错：", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		// 处理响应
		// 这里你可以解析 resp.Body 中的响应数据或者处理其他逻辑
		// 例如，检查响应状态码和返回的 JSON 数据等

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "请求失败"})
			return
		}

		c.String(http.StatusOK, "success")
	})

	//被动消息回复
	r.POST("/spark/answer", func(c *gin.Context) {
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
			//获取星火AI
			answer, _ := util.SparkAnswer(reqMsg.Content)
			c.JSON(http.StatusOK, gin.H{
				"ToUserName":   reqMsg.FromUserName,
				"FromUserName": reqMsg.ToUserName,
				"CreateTime":   reqMsg.CreateTime,
				"MsgType":      "text",
				"Content":      answer,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"ToUserName":   reqMsg.FromUserName,
				"FromUserName": reqMsg.ToUserName,
				"CreateTime":   reqMsg.CreateTime,
				"MsgType":      "text",
				"Content":      "暂时不支持文本之外的回复",
			})
		}
	})

	r.Run(":80")
}
