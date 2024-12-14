package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ApiServer interface {
	ListenAndServe() error
}

func InitServer(port string, router *gin.Engine) ApiServer {
	return &http.Server{
		Addr:           port,
		Handler:        router,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
