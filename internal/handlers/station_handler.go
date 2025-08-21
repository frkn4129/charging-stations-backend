package handlers

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "charging-stations-backend/internal/services"
    "charging-stations-backend/internal/models"
    "googlemaps.github.io/maps"
    "fmt"
    "log"
)

type StationHandler struct {
    stationService *services.StationService
    mapService     *services.MapService
    reviewHandler  *ReviewHandler
}

func NewStationHandler(ss *services.StationService, ms *services.MapService, rh *ReviewHandler) *StationHandler {
    return &StationHandler{
        stationService: ss,
        mapService:     ms,
        reviewHandler:  rh,
    }
}

func (h *StationHandler) GetStations(c *gin.Context) {
    // Tüm istasyonları getir, ama istatistikler olmadan
    stations := h.stationService.GetStations()
	fmt.Println("stations", stations)
    c.JSON(http.StatusOK, stations)
}

func (h *StationHandler) GetStationDetails(c *gin.Context) {
    stationID := c.Param("id")
    
    // İstasyonu bul
    station := h.stationService.GetStation(stationID)
    if station == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "İstasyon bulunamadı"})
        return
    }

    // İstasyon için güncel istatistikleri al
    stats, err := h.reviewHandler.GetStationStats(stationID)
    if err != nil {
        log.Printf("Error getting stats for station %s: %v", stationID, err)
        stats = &models.StationStats{
            AverageRating: 0,
            ReviewCount:   0,
        }
    }

    // İstatistikleri istasyon bilgilerine ekle
    station.AverageRating = stats.AverageRating
    station.ReviewCount = stats.ReviewCount

    // Debug için yanıtı logla
    log.Printf("Station details response: %+v", station)

    c.JSON(http.StatusOK, station)
}

func (h *StationHandler) GetNearbyStations(c *gin.Context) {
    lat, err := strconv.ParseFloat(c.Query("lat"), 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
        return
    }

    lon, err := strconv.ParseFloat(c.Query("lng"), 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
        return
    }

    limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
        return
    }

    stations := h.stationService.GetNearbyStations(lat, lon, limit)
    c.JSON(http.StatusOK, stations)
}

func (h *StationHandler) GetRoute(c *gin.Context) {
    var req models.RouteRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    route, err := h.mapService.GetRoute(
        maps.LatLng{Lat: req.OriginLat, Lng: req.OriginLng},
        maps.LatLng{Lat: req.DestinationLat, Lng: req.DestinationLng},
    )

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, route)
}

func (h *StationHandler) CalculateDistance(c *gin.Context) {
    var req struct {
        Lat1 float64 `json:"lat1" binding:"required"`
        Lon1 float64 `json:"lon1" binding:"required"`
        Lat2 float64 `json:"lat2" binding:"required"`
        Lon2 float64 `json:"lon2" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        log.Printf("Invalid request data: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "error": fmt.Sprintf("Invalid request data: %v", err),
        })
        return
    }

    log.Printf("Received distance calculation request: %+v", req)

    result, err := h.mapService.GetDistance(
        maps.LatLng{Lat: req.Lat1, Lng: req.Lon1},
        maps.LatLng{Lat: req.Lat2, Lng: req.Lon2},
    )

    if err != nil {
        log.Printf("Distance calculation failed: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": fmt.Sprintf("Distance calculation failed: %v", err),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "distance": result.Distance,
        "duration": result.Duration,
    })
} 