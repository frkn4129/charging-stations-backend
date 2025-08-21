package handlers

import (
    "database/sql"
    "github.com/gin-gonic/gin"
    "charging-stations-backend/internal/models"
    "log"
    "net/http"
)

type ReviewHandler struct {
    db *sql.DB
}

func NewReviewHandler(db *sql.DB) *ReviewHandler {
    return &ReviewHandler{
        db: db,
    }
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
    stationID := c.Param("id")
    if stationID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Station ID is required"})
        return
    }

    var review struct {
        Rating  float64 `json:"rating" binding:"required,min=1,max=5"`
        Comment string  `json:"comment"`
    }

    if err := c.ShouldBindJSON(&review); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Review oluşturma işlemi...
    query := `
        INSERT INTO reviews (station_id, rating, comment, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
        RETURNING id
    `

    var id int
    err := h.db.QueryRow(query, stationID, review.Rating, review.Comment).Scan(&id)
    if err != nil {
        log.Printf("Error creating review: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
        return
    }

    // Güncel istatistikleri al
    stats, err := h.GetStationStats(stationID)
    if err != nil {
        log.Printf("Error getting updated stats: %v", err)
    }

    // Başarılı yanıtla birlikte güncel istatistikleri de gönder
    c.JSON(http.StatusCreated, gin.H{
        "id": id,
        "message": "Review created successfully",
        "stats": stats,
    })
}

func (h *ReviewHandler) GetStationReviews(c *gin.Context) {
    stationID := c.Param("id")
    log.Printf("İstasyon yorumları istendi: StationID=%s", stationID)

    query := `
        SELECT id, station_id, rating, comment, created_at, updated_at
        FROM reviews
        WHERE station_id = $1
        ORDER BY created_at DESC`

    rows, err := h.db.Query(query, stationID)
    if err != nil {
        log.Printf("PostgreSQL okuma hatası: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Yorumlar alınamadı",
            "details": err.Error(),
        })
        return
    }
    defer rows.Close()

    var reviews []models.Review
    for rows.Next() {
        var review models.Review
        err := rows.Scan(
            &review.ID,
            &review.StationID,
            &review.Rating,
            &review.Comment,
            &review.CreatedAt,
            &review.UpdatedAt,
        )
        if err != nil {
            log.Printf("Row okuma hatası: %v", err)
            continue
        }
        reviews = append(reviews, review)
    }

    log.Printf("%d yorum bulundu", len(reviews))
    c.JSON(http.StatusOK, reviews)
}

func (h *ReviewHandler) GetStationStats(stationID string) (*models.StationStats, error) {
    query := `
        SELECT COALESCE(ROUND(AVG(rating)::numeric, 1), 0) as average_rating, 
               COUNT(*) as review_count
        FROM reviews
        WHERE station_id = $1`

    stats := &models.StationStats{
        AverageRating: 0,
        ReviewCount:   0,
    }

    err := h.db.QueryRow(query, stationID).Scan(&stats.AverageRating, &stats.ReviewCount)
    if err != nil {
        return nil, err
    }

    return stats, nil
} 