package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

/**
*  WebAPI 接口调用示例 接口文档（必看）：https://www.xfyun.cn/doc/spark/Web.html
* 错误码链接：https://www.xfyun.cn/doc/spark/%E6%8E%A5%E5%8F%A3%E8%AF%B4%E6%98%8E.html（code返回错误码时必看）
* @author iflytek
 */

var (
	hostUrl   = "ws://spark-api.xf-yun.com/v2.1/chat"
	appid     = "6015d7f0"
	apiSecret = "ZWQ3ZTVmZGFmZmE4ZTIzZGJjNzJlN2Q3"
	apiKey    = "45ec366167ba190909c776cc42a13484"
)

func SparkAnswer(question string) (string, error) {
	d := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}

	// 建立 WebSocket 连接
	conn, resp, err := d.Dial(assembleAuthUrl(hostUrl, apiKey, apiSecret), nil)
	if err != nil {
		return "", fmt.Errorf("WebSocket 连接失败: %s", err.Error())
	}
	defer conn.Close()

	if resp.StatusCode != 101 {
		return "", fmt.Errorf("WebSocket 连接失败，状态码：%d", resp.StatusCode)
	}

	// 使用通道来传递消息和错误
	msgCh := make(chan string)
	errCh := make(chan error)

	// 发送消息的协程
	go func() {
		data := genParams(appid, question)
		if err := conn.WriteJSON(data); err != nil {
			errCh <- fmt.Errorf("发送数据失败: %s", err.Error())
			return
		}
	}()

	// 接收消息的协程
	go func() {
		var answer string
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				errCh <- fmt.Errorf("读取消息失败: %s", err.Error())
				return
			}

			var data map[string]interface{}
			err = json.Unmarshal(msg, &data)
			if err != nil {
				errCh <- fmt.Errorf("解析 JSON 失败: %s", err.Error())
				return
			}

			payload := data["payload"].(map[string]interface{})
			choices := payload["choices"].(map[string]interface{})
			header := data["header"].(map[string]interface{})
			code := header["code"].(float64)

			if code != 0 {
				errCh <- fmt.Errorf("错误代码：%f", code)
				return
			}

			status := choices["status"].(float64)
			text := choices["text"].([]interface{})
			content := text[0].(map[string]interface{})["content"].(string)

			if status != 2 {
				answer += content
			} else {
				answer += content
				usage := payload["usage"].(map[string]interface{})
				temp := usage["text"].(map[string]interface{})
				totalTokens := temp["total_tokens"].(float64)
				fmt.Println("total_tokens:", totalTokens)
				msgCh <- answer
				return
			}
		}
	}()

	// 监听消息或错误
	select {
	case answer := <-msgCh:
		fmt.Println("spark的回答：", answer)
		return answer, nil
	case err := <-errCh:
		return "", err
	}
}

// 生成参数
func genParams(appid, question string) map[string]interface{} { // 根据实际情况修改返回的数据结构和字段名

	messages := []Message{
		{Role: "user", Content: question},
	}

	data := map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
		"header": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"app_id": appid, // 根据实际情况修改返回的数据结构和字段名
		},
		"parameter": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"chat": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
				"domain":      "generalv2",  // 根据实际情况修改返回的数据结构和字段名
				"temperature": float64(0.8), // 根据实际情况修改返回的数据结构和字段名
				"top_k":       int64(6),     // 根据实际情况修改返回的数据结构和字段名
				"max_tokens":  int64(2048),  // 根据实际情况修改返回的数据结构和字段名
				"auditing":    "default",    // 根据实际情况修改返回的数据结构和字段名
			},
		},
		"payload": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"message": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
				"text": messages, // 根据实际情况修改返回的数据结构和字段名
			},
		},
	}
	return data // 根据实际情况修改返回的数据结构和字段名
}

// 创建鉴权url  apikey 即 hmac username
func assembleAuthUrl(hosturl string, apiKey, apiSecret string) string {
	ul, err := url.Parse(hosturl)
	if err != nil {
		fmt.Println(err)
	}
	//签名时间
	date := time.Now().UTC().Format(time.RFC1123)
	//date = "Tue, 28 May 2019 09:10:42 MST"
	//参与签名的字段 host ,date, request-line
	signString := []string{"host: " + ul.Host, "date: " + date, "GET " + ul.Path + " HTTP/1.1"}
	//拼接签名字符串
	sgin := strings.Join(signString, "\n")
	// fmt.Println(sgin)
	//签名结果
	sha := HmacWithShaTobase64("hmac-sha256", sgin, apiSecret)
	// fmt.Println(sha)
	//构建请求参数 此时不需要urlencoding
	authUrl := fmt.Sprintf("hmac username=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"", apiKey,
		"hmac-sha256", "host date request-line", sha)
	//将请求参数使用base64编码
	authorization := base64.StdEncoding.EncodeToString([]byte(authUrl))

	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	//将编码后的字符串url encode后添加到url后面
	callurl := hosturl + "?" + v.Encode()
	fmt.Println("callurl是：", callurl)
	return callurl
}

func HmacWithShaTobase64(algorithm, data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

func ReadResp(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("code=%d,body=%s", resp.StatusCode, string(b))
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
