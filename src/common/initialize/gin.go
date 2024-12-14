package initialize

import (
	"bfv-bot/common/global"
	"bfv-bot/common/router"
	"bfv-bot/common/utils"
	"bfv-bot/model/common/resp"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

func Routers() *gin.Engine {
	gin.SetMode(global.GConfig.Server.GinMode)
	ginDefault := gin.New()
	ginDefault.Use(GinLogger(), GinRecovery(true))
	ginDefault.Use(utils.Options())

	prefix := ginDefault.Group("api")
	publicGroup := prefix.Group("")
	router.RouterGroupApp.EventRouter.InitPublicRouter(publicGroup)

	global.GLog.Info("路由注册完成")
	return ginDefault
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		global.GLog.Debug(fmt.Sprintf("msg: %s\n", string(body)))

		// 将请求体重新设置到请求中
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		cost := time.Since(start).Milliseconds()

		if query != "" {
			path = path + "?" + query
		}

		// 写入日志
		global.GLog.Info(fmt.Sprintf("[GIN] |%s|%9s |%17s |%s \"%s\"", colorfulStatus(c.Writer.Status()),
			fmt.Sprintf("%dms", cost),
			c.ClientIP(), colorfulMethod(c.Request.Method), path))
	}
}

func colorfulStatus(status int) string {
	var colorCode string
	if status == 200 {
		colorCode = "\033[0;42m %d \033[0m"
	} else if status == 404 {
		colorCode = "\033[0;43;30m %d \033[0m"
	} else {
		colorCode = "\033[0;41m %d \033[0m"
	}
	return fmt.Sprintf(colorCode, status)
}

func colorfulMethod(method string) string {
	var colorCode string
	switch method {
	case "POST":
		colorCode = "\033[0;48;2;120;220;232;38;2;255;255;255m %-8s\033[0m"
	case "GET":
		colorCode = "\033[0;48;2;252;152;103;38;2;255;255;255m %-8s\033[0m"
	case "PUT":
		colorCode = "\033[0;43;30m %-8s\033[0m"
	case "DELETE":
		colorCode = "\033[0;48;2;255;97;136;97m %-8s\033[0m"
	case "OPTIONS":
		colorCode = "\033[0;48;2;147;146;147;38;2;98;88;98m %-8s\033[0m"
	default:
		colorCode = " %-8s"
	}
	return fmt.Sprintf(colorCode, method)
}

// GinRecovery recover掉项目可能出现的panic
// 此函数是捕获panic，根据gin框架内的Recovery修改的
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequestByteArr, _ := httputil.DumpRequest(c.Request, true)
				httpRequestString := strings.ReplaceAll(string(httpRequestByteArr), "\r\n", " ")

				var errInfo string
				if brokenPipe {
					errInfo = c.Request.URL.Path
				} else {
					errInfo = "[Recovery from panic]"
				}

				global.GLog.Error(errInfo, zap.Any("error", err), zap.String("request", httpRequestString))

				if brokenPipe {
					_ = c.Error(err.(error))
					c.Abort()
					return
				}

				if stack {
					stack := string(debug.Stack())
					global.GLog.Error(stack)
				}

				c.JSON(http.StatusInternalServerError, resp.EmptyResponse{})
			}
		}()
		c.Next()
	}
}
