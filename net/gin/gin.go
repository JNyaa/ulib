package gin_utils

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/3JoB/telebot/pkg"
	"github.com/3JoB/ulib/json"
)

func Bind(c *gin.Context, v any) {
	data, _ := io.ReadAll(c.Request.Body)
	json.Unmarshal(data, v)
}

func Body(c *gin.Context) string {
	body := make([]byte, c.Request.ContentLength)
	c.Request.Body.Read(body)
	return pkg.String(body)
}