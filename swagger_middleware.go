package gin_swagger

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Config holds all options for swagger protection & route.
type Config struct {
	// If false, the swagger route will always return 403.
	Enabled bool

	// Path for swagger route: e.g. "/swagger/*any"
	// If empty, default is "/swagger/*any".
	Path string

	// Optional list of allowed client IPs.
	AllowedIPs []string

	// Optional static token for header-based auth:
	// Client must send:  X-API-TOKEN: <AuthToken>
	AuthToken string
}

// Option function for functional options pattern.
type Option func(*Config)

// WithPath sets swagger route path.
func WithPath(path string) Option {
	return func(c *Config) {
		c.Path = path
	}
}

// WithAllowedIPs sets list of allowed IPs.
func WithAllowedIPs(ips ...string) Option {
	return func(c *Config) {
		c.AllowedIPs = ips
	}
}

// WithAuthToken sets a static auth token.
func WithAuthToken(token string) Option {
	return func(c *Config) {
		c.AuthToken = token
	}
}

// WithEnabled explicitly sets Enabled.
func WithEnabled(enabled bool) Option {
	return func(c *Config) {
		c.Enabled = enabled
	}
}

// defaultConfig makes a baseline config.
func defaultConfig() *Config {
	return &Config{
		Enabled: true,
		Path:    "/swagger/*any",
	}
}

// NewMiddleware creates a gin.HandlerFunc from a Config.
// You can use this directly if you want to plug it yourself.
func NewMiddleware(cfg *Config) gin.HandlerFunc {
	if cfg == nil {
		cfg = defaultConfig()
	}
	if cfg.Path == "" {
		cfg.Path = "/swagger/*any"
	}

	return func(c *gin.Context) {
		// 1. Disabled?
		if !cfg.Enabled {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Swagger UI is disabled",
			})
			return
		}

		// 2. IP allow list
		if len(cfg.AllowedIPs) > 0 {
			clientIP := c.ClientIP()
			allowed := false
			for _, ip := range cfg.AllowedIPs {
				if clientIP == ip {
					allowed = true
					break
				}
			}
			if !allowed {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "IP not allowed",
				})
				return
			}
		}

		// 3. Token check
		if cfg.AuthToken != "" {
			token := c.GetHeader("X-API-TOKEN")
			if token != cfg.AuthToken {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Unauthorized swagger access",
				})
				return
			}
		}

		c.Next()
	}
}

// AttachSwagger is the GENERIC function you reuse in all projects.
//
// - router: any gin.IRouter (Engine, Group, etc.)
// - handler: any gin.HandlerFunc (gin-swagger, Redoc, your own handler…)
// - opts: functional options to configure behavior and path.
func AttachSwagger(router gin.IRouter, handler gin.HandlerFunc, opts ...Option) {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	mw := NewMiddleware(cfg)

	// NOTE: Path must be like "/swagger/*any" for gin-swagger style handlers.
	router.GET(cfg.Path, mw, handler)
}
