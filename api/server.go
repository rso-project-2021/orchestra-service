package api

import (
	"orchestra-service/config"
	"orchestra-service/proto"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	config     config.Config
	router     *gin.Engine
	grpcClient proto.LocationServiceClient
}

func NewServer(config config.Config, grpcClient proto.LocationServiceClient) (*Server, error) {

	gin.SetMode(config.GinMode)
	router := gin.Default()

	server := &Server{
		config:     config,
		grpcClient: grpcClient,
	}

	// Setup routing for server.
	v1 := router.Group("v1")
	{
		v1.GET("/quickreserve", server.QuickReserve)
	}

	// Setup health check routes.
	health := router.Group("health")
	{
		health.GET("/live", server.Live)
		health.GET("/ready", server.Ready)
	}

	// Setup metrics routes.
	metrics := router.Group("metrics")
	{
		metrics.GET("/", func(ctx *gin.Context) {
			handler := promhttp.Handler()
			handler.ServeHTTP(ctx.Writer, ctx.Request)
		})
	}

	server.router = router
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
