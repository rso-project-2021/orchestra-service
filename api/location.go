package api

import (
	"log"
	"net/http"
	"orchestra-service/proto"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

type quickReserveRequest struct {
	ID  int64   `json:"id" binding:"required"`
	Lat float32 `json:"lat" binding:"required"`
	Lng float32 `json:"lng" binding:"required"`
}

func (server *Server) QuickReserve(ctx *gin.Context) {

	// Check if request has all required fields in json body.
	var req quickReserveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("Sem tu 1")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		ctx.Abort()
		return
	}

	myLocation := proto.Location{
		Id:  req.ID,
		Lat: req.Lat,
		Lng: req.Lng,
	}

	locationRequest := proto.LocationRequest{
		MyLocation: &myLocation,
		Locations:  []*proto.Location{},
	}

	closest, err := server.grpcClient.FindClosest(context.Background(), &locationRequest)
	if err != nil {
		log.Println("Sem tu 2")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, closest)
}
