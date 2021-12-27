package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) Live(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
}

func (server *Server) Ready(ctx *gin.Context) {

	// Check connection to station service.
	url := fmt.Sprintf("%s/health/ready", server.config.StationsAddress)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": "DOWN"})
		ctx.Abort()
		return
	}
	defer resp.Body.Close()

	// Check connection to reservation service.
	url = fmt.Sprintf("%s/health/ready", server.config.ReservationsAddress)
	resp, err = http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": "DOWN"})
		ctx.Abort()
		return
	}
	defer resp.Body.Close()

	ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
}
