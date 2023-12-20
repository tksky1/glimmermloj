package controller

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func Submit(c *gin.Context) {

	token, _ := c.GetQuery("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "token not found"})
		return
	}
	ch := make(chan interface{})
	go getUserinfo(ch, token)
	username1 := <-ch
	username := username1.(string)
	userID1 := <-ch
	userID := userID1.(int64)
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "token invalid"})
		return
	}
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "cant read body"})
		return
	}
	go Judge(userID, username, buf)

	// 返回成功的响应
	c.JSON(http.StatusOK, gin.H{"message": "submit success"})
}
