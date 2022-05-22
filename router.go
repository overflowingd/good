package good

import (
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	if _, debug := os.LookupEnv("DEBUG"); !debug {
		gin.SetMode(gin.ReleaseMode)
	}
}

func NewRouter() (*gin.Engine, error) {
	return gin.Default(), nil
}
