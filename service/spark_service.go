package service

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go_spark_ai/util"

	"github.com/gin-gonic/gin"
)

// JsonResult 返回结构
type JsonResult struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// RequestBody 微信回复文本消息结构体
type RequestBody struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string
	MsgId        string
}

// ResponseBody 微信回复文本消息结构体
type ResponseBody struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string
}

// IndexHandler 讯飞星火接口
func SparkAIHandler(c *gin.Context) {
	var reqMsg RequestBody
	err := c.ShouldBindJSON(&reqMsg)
	if err != nil {
		log.Printf("[消息接收] - JSON数据包解析失败: %v\n", err)
		return
	}

	log.Printf("[消息接收] - 收到消息, 消息类型为: %s, 消息内容为: %s\n", reqMsg.MsgType, reqMsg.Content)

	// 获取星火AI答案
	answer := util.SparkAnswer(reqMsg.Content)

	// 对接收的消息进行被动回复
	WXMsgReply(c, reqMsg.FromUserName, reqMsg.ToUserName, answer)

	c.String(http.StatusOK, "success") // 发送成功响应
}

// WXMsgReply 微信消息回复
func WXMsgReply(c *gin.Context, fromUser, toUser, content string) {
	respMsg := ResponseBody{
		ToUserName:   toUser,
		FromUserName: fromUser,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      content,
	}

	msg, err := json.Marshal(&respMsg)
	if err != nil {
		log.Printf("[消息回复] - 将对象进行JSON编码出错: %v\n", err)
		return
	}
	_, _ = c.Writer.Write(msg)
}
