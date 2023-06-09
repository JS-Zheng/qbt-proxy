package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

type Server struct {
	router *gin.Engine
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (server *Server) Run(cfg Config) {
	if cfg.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	// Add CORS middleware to allow CORS requests
	router.Use(createCorsMw(cfg))
	server.router = router

	// Set up the API proxy routes.
	setupRoutes(router, cfg)

	if shouldStartHTTPSServer(cfg) {
		startHTTPSServer(router, cfg)
	} else {
		startHTTPServer(router, cfg)
	}
}

func createCorsMw(cfg Config) gin.HandlerFunc {
	corsCfg := cors.DefaultConfig()

	if len(cfg.CORSAllowOrigins) > 0 {
		if contains(cfg.CORSAllowOrigins, "*") {
			corsCfg.AllowAllOrigins = true
		} else {
			corsCfg.AllowOrigins = cfg.CORSAllowOrigins
		}
	} else {
        corsCfg.AllowAllOrigins = true
    }

	if len(cfg.CORSAllowMethods) > 0 {
		corsCfg.AllowMethods = cfg.CORSAllowMethods
	}

	if len(cfg.CORSAllowHeaders) > 0 {
		corsCfg.AllowHeaders = cfg.CORSAllowHeaders
	}

	if len(cfg.CORSExposeHeaders) > 0 {
		corsCfg.ExposeHeaders = cfg.CORSExposeHeaders
	}

	corsCfg.AllowCredentials = cfg.CORSAllowCredentials

    log.Println("CORS Config:")
	log.Printf("%+v\n", corsCfg)

	return cors.New(corsCfg)
}

// Configures the proxy routes for the given router and configuration.
func setupRoutes(router *gin.Engine, cfg Config) {
	router.Any("/api/v2/*apiPath", func(c *gin.Context) {
		apiPath := c.Param("apiPath")
		if apiPath == "/auth/login" && c.Request.Method == "POST" {
			LoginProxy(cfg)(c)
			return
		}

		// Route all other requests to the ApiProxy middleware function.
		ApiProxy(cfg)(c)
	})
}

func shouldStartHTTPSServer(cfg Config) bool {
	return cfg.HTTPSPort > 0 && cfg.TLSCert != "" && cfg.TLSKey != ""
}

func startHTTPSServer(router *gin.Engine, cfg Config) {
	httpsPortStr := cfg.GetHTTPSPortStr()
	log.Printf("HTTPS server starting on port %s", httpsPortStr)
	err := router.RunTLS(httpsPortStr, cfg.TLSCert, cfg.TLSKey)
	if err != nil {
		log.Fatalf("Failed to start HTTPS server: %v", err)
	}
}

func startHTTPServer(router *gin.Engine, cfg Config) {
	httpPortStr := cfg.GetHTTPPortStr()
	log.Printf("HTTP server starting on port %s", httpPortStr)
	err := router.Run(httpPortStr)
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func HandleError(c *gin.Context, code int, msg string, err error) {
	log.Println(err)
	c.JSON(code, ErrorResponse{
		Error: msg,
	})
}

func contains(slice []string, target string) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}
