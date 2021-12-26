package api

import (
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
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		ctx.Abort()
		return
	}

	origin := proto.Location{
		Id:  req.ID,
		Lat: req.Lat,
		Lng: req.Lng,
	}

	destinations := []*proto.Location{
		{ // Supernova
			Id:  1,
			Lat: 46.0364004,
			Lng: 14.5252039,
		},
		{ // Zdravstveni dom Vič
			Id:  2,
			Lat: 46.0463459,
			Lng: 14.4818508,
		},
		{ // Logatec
			Id:  3,
			Lat: 45.918809,
			Lng: 14.2161334,
		},
	}

	locationRequest := proto.LocationRequest{
		Origin:       &origin,
		Destinations: destinations,
	}

	closest, err := server.grpcClient.FindClosest(context.Background(), &locationRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, closest)
}
