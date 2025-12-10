package transport

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/noctisnyxx/auth-service/configs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(handler *UserHandler) *gin.Engine {
	gin.SetMode(configs.GIN_MODE)
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Range"}
	config.ExposeHeaders = []string{"Content-Length", "Content-Range", "Content-Disposition"}

	r.Use(cors.New(config))
	r.Use(gin.Recovery())
	r.Use(ErrorLoggingMiddleware())
	api := r.Group("/api")
	userGroup := api.Group("/user")
	{
		userGroup.GET("", handler.FindUsers)
		userGroup.POST("/register", handler.RegisterUser)
	}

	//swagger
	fmt.Println("swagger URL: http://localhost:" + configs.BE_PORT + "/api/swagger/index.html")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.NewHandler()))

	return r
}

func ErrorLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		status := c.Writer.Status()
		if status >= 400 {
			errs := c.Errors.ByType(gin.ErrorTypeAny)

			if len(errs) > 0 {
				log.Printf("[ERROR] %d %s | %s", status, c.Request.Method, c.Request.URL.Path)
				for _, e := range errs {
					log.Printf("  → %v", e.Err)
				}
			} else {
				log.Printf("[ERROR] %d %s | %s", status, c.Request.Method, c.Request.URL.Path)
			}
		}
	}
}
