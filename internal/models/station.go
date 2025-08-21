package models

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
    AverageRating          float64 `json:"average_rating"`
    ReviewCount            int     `json:"review_count"`
}

type StationStats struct {
    AverageRating float64 `json:"average_rating"`
    ReviewCount   int     `json:"review_count"`
}

type TrugoResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
    Data    struct {
        Stations []struct {
            Brand                    string  `json:"brand"`
            ID                      int     `json:"id"`
            StationID              string  `json:"station_id"`
            Name                    string  `json:"name"`
            Latitude               float64 `json:"latitude"`
            Longitude              float64 `json:"longitude"`
            Status                 string  `json:"status"`
            ConnectorList          string  `json:"connector_list"`
            ErrorDeviceCount       int     `json:"error_device_count"`
            ACAvailableSocketsCount int    `json:"ac_available_sockets_count"`
            DCAvailableSocketsCount int    `json:"dc_available_sockets_count"`
            UnavailableDeviceCount  int    `json:"unavailable_device_count"`
            TotalConnectorsCount    int    `json:"total_connectors_count"`
            StationColor           string  `json:"station_color"`
        } `json:"stations"`
    } `json:"data"`
}

type RouteRequest struct {
    OriginLat      float64 `json:"origin_lat"`
    OriginLng      float64 `json:"origin_lng"`
    DestinationLat float64 `json:"destination_lat"`
    DestinationLng float64 `json:"destination_lng"`
} 