package controller

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func getUserinfo(ch chan interface{}, token string) {

	var userID int64
	userID = 0
	nickName := ""

	url := "https://lab.glimmer.org.cn:9521/user/userInfo"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Authorization", token)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: 2 * time.Second, Transport: tr}
	resp, err := client.Do(req)

	if resp == nil || resp.Body == nil {
		println("请求光点平台失败")
		ch <- ""
		ch <- 0
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	// 解析JSON响应
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		println("光点平台返回无法解析")
	}

	if result["code"] == nil {
		println("光点平台返回无法解析")
		fmt.Printf("%+v", result)
		ch <- ""
		ch <- 0
		return
	}

	// 检查响应状态码
	code := result["code"].(float64)
	if code != 200 {
		println("用户token错误")
		ch <- ""
		ch <- 0
		return
	}

	// 解析用户信息
	testData := result["data"]
	if testData == nil {
		println("光点平台返回无法解析")
		fmt.Printf("%+v", result)
		ch <- ""
		ch <- 0
		return
	}
	fmt.Printf("%+v", result)
	data := result["data"].(map[string]interface{})
	userIDfloat := data["userId"].(float64)
	userID = int64(userIDfloat)
	nickName = data["nickName"].(string)

	ch <- nickName
	ch <- userID
}
