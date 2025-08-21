package models

import (
    "time"
)

type Review struct {
    ID        int       `json:"id"`
    StationID string    `json:"station_id"`
    Rating    float64   `json:"rating"`
    Comment   string    `json:"comment"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
} 