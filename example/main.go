package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	gin_swagger "github.com/raza001/gin_swaggers"
	"github.com/raza001/gin_swaggers/example/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Hello godoc
// @Summary Say hello
// @Description Returns hello message
// @Tags example
// @Produce json
// @Success 200 {object} map[string]string
// @Router /hello [get]
func Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "hello world"})
}

func main() {
	r := gin.Default()

	// If you're behind a reverse proxy in production, configure trusted proxies:
	// r.SetTrustedProxies([]string{"127.0.0.1", "172.17.0.1"}) // example

	// Set docs metadata (if you run swag init it will override with generated docs)
	docs.SwaggerInfo.Title = "My API"
	docs.SwaggerInfo.Version = "1.0"

	// Public example endpoint
	r.GET("/hello", Hello)

	// Attach swagger: protect UI but leave doc.json public (so the UI can fetch it)
	gin_swagger.AttachSwagger(r,
		// options:
		gin_swagger.WithPath("/swagger/*any"),
		gin_swagger.WithAllowedIPs("127.0.0.1", "::1"),
		// gin_swagger.WithAuthToken("my-secret-token"), // optional header token
		gin_swagger.WithProtectDocJSON(false), // false => doc.json public
	)

	// Additionally, you can mount raw gin-swagger handler directly on a route
	// to test without protection. Example:
	r.GET("/swagger-unprotected/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
