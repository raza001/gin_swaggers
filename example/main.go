package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/you/gin-swagger-protect/gin_swagger" // change to your module path
)

// A minimal handler that stands in for the swagger handler (replace with ginSwagger.WrapHandler)
func fakeSwaggerHandler(c *gin.Context) {
	c.String(http.StatusOK, "swagger content for %s", c.Request.RequestURI)
}

func main() {
	r := gin.Default()

	// Attach to top-level engine
	gin_swagger.AttachSwagger(r, fakeSwaggerHandler,
		gin_swagger.WithPath("/swagger/*any"),
		gin_swagger.WithAllowedIPs("127.0.0.1", "192.168.0.0/16"),
		gin_swagger.WithAuthToken("my-secret-token"),
	)

	// Or attach to a group (e.g., /admin)
	admin := r.Group("/admin")
	gin_swagger.AttachSwagger(admin, fakeSwaggerHandler,
		gin_swagger.WithPath("/swagger/*any"),
		gin_swagger.WithEnabled(true),
		gin_swagger.WithAuthToken("admin-secret"),
	)

	_ = r.Run(":8080")
}
