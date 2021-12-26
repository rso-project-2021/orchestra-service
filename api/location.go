package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"orchestra-service/proto"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
)

type quickReserveRequest struct {
	ID  int64   `json:"id" binding:"required"`
	Lat float32 `json:"lat" binding:"required"`
	Lng float32 `json:"lng" binding:"required"`
}

type LocationItem struct {
	ID  int64   `json:"station_id" binding:"required"`
	Lat float32 `json:"lat" binding:"required"`
	Lng float32 `json:"lng" binding:"required"`
}

type LocationList struct {
	List []*LocationItem `json:"list" binding:"required"`
}

type GraphResponse struct {
	Data LocationList `json:"data" binding:"required"`
}

func (server *Server) QuickReserve(ctx *gin.Context) {

	// Check if request has all required fields in json body.
	var req quickReserveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		ctx.Abort()
		return
	}

	// Get nearby stations.
	graphQuery, _ := json.Marshal(map[string]string{
		"query": `
			{
				list {
					station_id
					lat
					lng
				}
			}
        `,
	})

	url := fmt.Sprintf("%s/v1/graphql?lat=%f&lng=%f&offset=%d&limit=%d", server.config.StationsAddress, req.Lat, req.Lng, 0, 5)
	resp, err := http.Post(url, "application/graphql", bytes.NewBuffer(graphQuery))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err})
		ctx.Abort()
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err})
		ctx.Abort()
		return
	}

	var res GraphResponse
	if err := json.Unmarshal(body, &res); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err})
		ctx.Abort()
		return
	}

	// Parse destinations into a slice.
	var destinations []*proto.Location
	for _, loc := range res.Data.List {
		destinations = append(destinations, &proto.Location{
			Id:  loc.ID,
			Lat: loc.Lat,
			Lng: loc.Lng,
		})
	}

	origin := proto.Location{
		Id:  req.ID,
		Lat: req.Lat,
		Lng: req.Lng,
	}

	locationRequest := proto.LocationRequest{
		Origin:       &origin,
		Destinations: destinations,
	}

	// Find station that is closest to origin.
	closest, err := server.grpcClient.FindClosest(context.Background(), &locationRequest)
	if err != nil {
		log.Warn("Google Maps service is not responding.")
		closest = QuickReserveFallback(&locationRequest)
		ctx.JSON(http.StatusOK, closest)
		return
	}

	ctx.JSON(http.StatusOK, closest)
}

func QuickReserveFallback(req *proto.LocationRequest) *proto.Location {
	return req.Destinations[0]
}
