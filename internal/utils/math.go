package utils

import "math"

func ToRadians(deg float64) float64 {
    return deg * math.Pi / 180
}

func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
    const R = 6371 // Dünya'nın yarıçapı (km)

    lat1Rad := ToRadians(lat1)
    lat2Rad := ToRadians(lat2)
    deltaLat := ToRadians(lat2 - lat1)
    deltaLon := ToRadians(lon2 - lon1)

    a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
        math.Cos(lat1Rad)*math.Cos(lat2Rad)*
        math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    return R * c
} 