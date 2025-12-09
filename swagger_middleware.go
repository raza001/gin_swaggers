package gin_swagger

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	// Importing docs is optional; the example main imports docs to set metadata.
	// If you generate docs with swag init, make sure module path matches.
)

// Config holds all options for swagger protection & route.
type Config struct {
	Enabled        bool     // If false, swagger route will always return 403.
	Path           string   // Path for swagger route: e.g. "/swagger/*any"
	AllowedIPs     []string // Optional list of allowed client IPs or CIDR (e.g. "192.168.0.0/16")
	AuthToken      string   // Optional static token header X-API-TOKEN
	ProtectDocJSON bool     // If true, /swagger/doc.json is protected by middleware as well
}

// Option function for functional options pattern.
type Option func(*Config)

// WithPath sets swagger route path.
func WithPath(path string) Option {
	return func(c *Config) {
		c.Path = path
	}
}

// WithAllowedIPs sets list of allowed IPs/CIDRs.
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

// WithProtectDocJSON sets whether the swagger JSON (doc.json) is protected.
func WithProtectDocJSON(protect bool) Option {
	return func(c *Config) {
		c.ProtectDocJSON = protect
	}
}

// defaultConfig makes a baseline config.
func defaultConfig() *Config {
	return &Config{
		Enabled:        true,
		Path:           "/swagger/*any",
		ProtectDocJSON: false,
	}
}

// parseAllowed converts AllowedIPs strings into two lists: individual IPs and CIDRs.
func parseAllowed(allowed []string) (ips []net.IP, cidrs []*net.IPNet) {
	for _, entry := range allowed {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		if strings.Contains(entry, "/") {
			if _, netw, err := net.ParseCIDR(entry); err == nil {
				cidrs = append(cidrs, netw)
				continue
			}
			// fallthrough to parse as single IP if CIDR parsing fails
		}
		if ip := net.ParseIP(entry); ip != nil {
			ips = append(ips, ip)
		}
	}
	return
}

// NewMiddleware creates a gin.HandlerFunc from a Config.
func NewMiddleware(cfg *Config) gin.HandlerFunc {
	if cfg == nil {
		cfg = defaultConfig()
	}
	if cfg.Path == "" {
		cfg.Path = "/swagger/*any"
	}

	ips, cidrs := parseAllowed(cfg.AllowedIPs)

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
			clientIPStr := c.ClientIP()
			clientIP := net.ParseIP(clientIPStr)
			allowed := false
			if clientIP != nil {
				for _, ip := range ips {
					if ip.Equal(clientIP) {
						allowed = true
						break
					}
				}
				if !allowed {
					for _, n := range cidrs {
						if n.Contains(clientIP) {
							allowed = true
							break
						}
					}
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

// AttachSwagger registers swagger UI & (optionally) doc.json with protection.
func AttachSwagger(router gin.IRouter, opts ...Option) {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	mw := NewMiddleware(cfg)
	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)

	// Register UI endpoints (GET/HEAD/OPTIONS)
	router.Handle("GET", cfg.Path, mw, func(c *gin.Context) {
		swaggerHandler(c)
	})
	router.Handle("HEAD", cfg.Path, mw, func(c *gin.Context) {
		swaggerHandler(c)
	})
	router.Handle("OPTIONS", cfg.Path, mw, func(c *gin.Context) {
		swaggerHandler(c)
	})

	// doc.json path used by gin-swagger is typically "/swagger/doc.json"
	// We construct it from the path prefix (cfg.Path without the "/*any")
	prefix := strings.TrimSuffix(cfg.Path, "/*any")
	if prefix == "" {
		prefix = "/"
	}
	docPath := prefix + "doc.json"

	if cfg.ProtectDocJSON {
		router.GET(docPath, mw, func(c *gin.Context) {
			swaggerHandler(c)
		})
	} else {
		// expose doc.json publicly so UI can fetch it even if UI is protected
		router.GET(docPath, func(c *gin.Context) {
			swaggerHandler(c)
		})
	}
}
