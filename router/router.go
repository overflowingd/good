package router

import "github.com/gin-gonic/gin"

type RouteBuilder func(*gin.Engine) error
