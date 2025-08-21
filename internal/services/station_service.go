package services

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "sort"
    "math"
    "strconv"
)

// API yanıt yapısı
type TrugoResponse struct {
    Status   string         `json:"status"`
    Message  string         `json:"message"`
    Data     StationsData   `json:"data"`
    Options  interface{}    `json:"options"`
}

// Stations verisi için ara yapı
type StationsData struct {
    Stations []Station `json:"stations"`
}

// İstasyon modeli
type Station struct {
    ID                      int     `json:"id"`
    StationID              string  `json:"station_id"`
    Name                   string  `json:"name"`
    Brand                  string  `json:"brand"`
    Latitude               float64 `json:"latitude"`
    Longitude              float64 `json:"longitude"`
    ConnectorList          string  `json:"connector_list"`
    ErrorDeviceCount       int     `json:"error_device_count"`
    ACAvailableSocketCount int     `json:"ac_available_sockets_count"`
    DCAvailableSocketCount int     `json:"dc_available_sockets_count"`
    UnavailableDeviceCount int     `json:"unavailable_device_count"`
    TotalConnectorsCount   int     `json:"total_connectors_count"`
    StationColor           string  `json:"station_color"`
    // Review için eklenen alanlar
    AverageRating          float64 `json:"average_rating,omitempty"`
    ReviewCount            int     `json:"review_count,omitempty"`
}

type StationService struct {
    stations []Station
    apiURL   string
}

func NewStationService() *StationService {
    return &StationService{
        apiURL: "https://emsp-api.trugo.com.tr/v1/public/csms/stations/fast/",
    }
}

func (s *StationService) loadStations() error {
    resp, err := http.Get(s.apiURL)
    if err != nil {
        return fmt.Errorf("API isteği hatası: %v", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("Response okuma hatası: %v", err)
    }

    var response TrugoResponse
    if err := json.Unmarshal(body, &response); err != nil {
        log.Printf("JSON parse hatası: %v", err)
        log.Printf("JSON içeriği: %s", string(body))
        return fmt.Errorf("JSON parse hatası: %v", err)
    }

    // İstasyonları stations alanından al
    s.stations = response.Data.Stations
    log.Printf("%d istasyon yüklendi", len(s.stations))
    return nil
}

func (s *StationService) GetStations() []Station {
    // İstasyonları yükle (cache mekanizması eklenebilir)
    if err := s.loadStations(); err != nil {
        log.Printf("İstasyonlar yüklenirken hata: %v", err)
        return []Station{}
    }
    return s.stations
}

func (s *StationService) GetStation(id string) *Station {
    // String ID'yi int'e çevir
    stationID, err := strconv.Atoi(id)
    if err != nil {
        log.Printf("Invalid station ID: %v", err)
        return nil
    }

    // İstasyonları yükle
    if err := s.loadStations(); err != nil {
        log.Printf("İstasyonlar yüklenirken hata: %v", err)
        return nil
    }

    for _, station := range s.stations {
        if station.ID == stationID {
            stationCopy := station
            return &stationCopy
        }
    }
    return nil
}

func (s *StationService) GetNearbyStations(lat, lon float64, limit int) []Station {
    // İstasyonları yükle
    if err := s.loadStations(); err != nil {
        log.Printf("İstasyonlar yüklenirken hata: %v", err)
        return []Station{}
    }

    // Mesafe hesaplama ve sıralama işlemleri...
    type stationDistance struct {
        station  Station
        distance float64
    }

    var stationsWithDistance []stationDistance
    for _, station := range s.stations {
        dist := calculateHaversineDistance(lat, lon, station.Latitude, station.Longitude)
        stationsWithDistance = append(stationsWithDistance, stationDistance{
            station:  station,
            distance: dist,
        })
    }

    // Mesafeye göre sırala
    sort.Slice(stationsWithDistance, func(i, j int) bool {
        return stationsWithDistance[i].distance < stationsWithDistance[j].distance
    })

    // İstenen sayıda istasyonu döndür
    var result []Station
    for i := 0; i < limit && i < len(stationsWithDistance); i++ {
        result = append(result, stationsWithDistance[i].station)
    }

    return result
}

// Haversine formülü ile iki nokta arası mesafe hesaplama
func calculateHaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
    const R = 6371 // Dünya'nın yarıçapı (km)

    dLat := (lat2 - lat1) * math.Pi / 180
    dLon := (lon2 - lon1) * math.Pi / 180

    a := math.Sin(dLat/2)*math.Sin(dLat/2) +
        math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
            math.Sin(dLon/2)*math.Sin(dLon/2)

    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    return R * c
} 