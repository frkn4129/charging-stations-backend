package services

import (
    "googlemaps.github.io/maps"
    "context"
    "charging-stations-backend/internal/utils"
    "os"
    "fmt"
    "log"
)

type MapService struct {
    client *maps.Client
}

type DistanceResult struct {
    Distance float64 `json:"distance"` // metre cinsinden
    Duration int64   `json:"duration"` // saniye cinsinden
}

func NewMapService() *MapService {
    apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
    if apiKey == "" {
        log.Fatal("Google Maps API key is missing")
    }

    client, err := maps.NewClient(maps.WithAPIKey(apiKey))
    if err != nil {
        log.Fatalf("Error creating Google Maps client: %v", err)
    }
    return &MapService{client: client}
}

func (s *MapService) CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
    return utils.CalculateDistance(lat1, lon1, lat2, lon2)
}

func (s *MapService) GetDistance(origin, destination maps.LatLng) (*DistanceResult, error) {
    // API çağrısı başarısız olursa kuş bakışı mesafeyi hesapla
    distance := utils.CalculateDistance(
        origin.Lat, 
        origin.Lng, 
        destination.Lat, 
        destination.Lng,
    )

    // Ortalama hız 60 km/s varsayarak süreyi hesapla
    duration := int64((distance / 60) * 3600) // saniye cinsinden

    result := &DistanceResult{
        Distance: distance,
        Duration: duration,
    }
    
    // Google Maps API'yi dene
    if s.client != nil {
        r := &maps.DistanceMatrixRequest{
            Origins:      []string{fmt.Sprintf("%f,%f", origin.Lat, origin.Lng)},
            Destinations: []string{fmt.Sprintf("%f,%f", destination.Lat, destination.Lng)},
            Mode:         maps.TravelModeDriving,
        }

        resp, err := s.client.DistanceMatrix(context.Background(), r)
        if err == nil && len(resp.Rows) > 0 && len(resp.Rows[0].Elements) > 0 {
            element := resp.Rows[0].Elements[0]
            if element.Status == "OK" {
                result.Distance = float64(element.Distance.Meters) / 1000
                result.Duration = int64(element.Duration.Seconds())
            }
        }
    }

    return result, nil
}

func (s *MapService) GetRoute(origin, destination maps.LatLng) ([]maps.LatLng, error) {
    r := &maps.DirectionsRequest{
        Origin:      origin.String(),
        Destination: destination.String(),
        Mode:        maps.TravelModeDriving,
    }

    resp, _, err := s.client.Directions(context.Background(), r)
    if err != nil {
        return nil, err
    }

    if len(resp) == 0 {
        return nil, fmt.Errorf("no route found")
    }

    var route []maps.LatLng
    for _, leg := range resp[0].Legs {
        for _, step := range leg.Steps {
            route = append(route, step.StartLocation)
            route = append(route, step.EndLocation)
        }
    }

    return route, nil
} 