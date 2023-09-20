package service

import (
	"encoding/json"
	"log"
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

// MsgBody 微信回复文本消息结构体
type MsgBody struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string
}

// IndexHandler 讯飞星火接口
func SparkAIHandler(c *gin.Context) {
	//res := &JsonResult{}

	var textMsg MsgBody
	err := c.ShouldBindJSON(&textMsg)
	if err != nil {
		log.Printf("[消息接收] - JSON数据包解析失败: %v\n", err)
		return
	}

	log.Printf("[消息接收] - 收到消息, 消息类型为: %s, 消息内容为: %s\n", textMsg.MsgType, textMsg.Content)

	// 获取星火AI答案
	answer := util.SparkAnswer(textMsg.Content)

	// 对接收的消息进行被动回复
	WXMsgReply(c, textMsg.FromUserName, textMsg.ToUserName, answer)

	return
}

// WXMsgReply 微信消息回复
func WXMsgReply(c *gin.Context, fromUser, toUser, content string) {
	repTextMsg := MsgBody{
		ToUserName:   toUser,
		FromUserName: fromUser,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      content,
	}

	msg, err := json.Marshal(&repTextMsg)
	if err != nil {
		log.Printf("[消息回复] - 将对象进行JSON编码出错: %v\n", err)
		return
	}
	_, _ = c.Writer.Write(msg)
}
