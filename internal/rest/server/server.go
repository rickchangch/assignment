package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewGinServer(port int) (*http.Server, *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	// TODO: Add CORS settings

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: engine,
	}

	return server, engine
}
