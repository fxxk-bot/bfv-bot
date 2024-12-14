package resp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type EmptyResponse struct {
}

type Data struct {
	File    string `json:"file"`
	Summary string `json:"summary"`
}

type ImageMessage struct {
	Type string `json:"type"`
	Data Data   `json:"data"`
}

type ImageReply struct {
	Reply []ImageMessage `json:"reply"`
}

type ReplyResponse struct {
	Reply string `json:"reply"`
}

func EmptyOk(c *gin.Context) {
	c.JSON(http.StatusOK, EmptyResponse{})
}

func ReplyOk(c *gin.Context, reply string) {
	c.JSON(http.StatusOK, ReplyResponse{reply})
}

func ReplyWithData(c *gin.Context, m map[string]interface{}) {
	c.JSON(http.StatusOK, m)
}

func ImageOk(c *gin.Context, path string, nickname string) {
	data := make([]ImageMessage, 1)
	data[0] = ImageMessage{Type: "image", Data: Data{
		File:    path,
		Summary: nickname,
	}}
	c.JSON(http.StatusOK, ImageReply{Reply: data})
}
