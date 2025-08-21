package main

import (
	"charging-stations-backend/internal/handlers"
	"charging-stations-backend/internal/middleware"
	"charging-stations-backend/internal/services"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL sürücüsü
	"golang.org/x/time/rate"
)

func createTables(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS reviews (
        id SERIAL PRIMARY KEY,
        station_id VARCHAR(255) NOT NULL,
        rating DECIMAL(3,1) NOT NULL,
        comment TEXT,
        created_at TIMESTAMP WITH TIME ZONE NOT NULL,
        updated_at TIMESTAMP WITH TIME ZONE NOT NULL
    );

    CREATE INDEX IF NOT EXISTS idx_reviews_station_id ON reviews(station_id);
    CREATE INDEX IF NOT EXISTS idx_reviews_created_at ON reviews(created_at DESC);
    `

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("tablo oluşturma hatası: %v", err)
	}

	log.Println("Veritabanı tabloları başarıyla oluşturuldu")
	return nil
}

func main() {
	// .env dosyasını yükle
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// API anahtarını kontrol et
	if os.Getenv("GOOGLE_MAPS_API_KEY") == "" {
		log.Fatal("GOOGLE_MAPS_API_KEY is required")
	}

	// PostgreSQL bağlantısı
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URI"))
	if err != nil {
		log.Fatal("PostgreSQL bağlantı hatası:", err)
	}
	defer db.Close()

	// Bağlantıyı test et
	if err := db.Ping(); err != nil {
		log.Fatal("PostgreSQL ping hatası:", err)
	}
	log.Println("PostgreSQL bağlantısı başarılı")

	// Tabloları oluştur
	if err := createTables(db); err != nil {
		log.Fatal(err)
	}

	// Debug log'ları aktif et
	gin.SetMode(gin.DebugMode)

	router := gin.Default()

	// Request logging middleware ekle
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	// CORS ayarları
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Services
	stationService := services.NewStationService()
	mapService := services.NewMapService()

	// Handlers
	reviewHandler := handlers.NewReviewHandler(db)
	stationHandler := handlers.NewStationHandler(stationService, mapService, reviewHandler)

	// Rate limiter oluştur (10 istek/dakika)
	rateLimiter := middleware.NewIPRateLimiter(rate.Every(time.Minute), 10)

	// Routes
	api := router.Group("/api")
	{
		api.GET("/stations", stationHandler.GetStations)
		api.GET("/stations/:id", stationHandler.GetStationDetails)
		api.GET("/stations/nearby", stationHandler.GetNearbyStations)

		// Distance ve route endpoint'lerine rate limit uygula
		distanceGroup := api.Group("/")
		distanceGroup.Use(middleware.RateLimitMiddleware(rateLimiter))
		{
			distanceGroup.POST("/stations/distance", stationHandler.CalculateDistance)
			distanceGroup.POST("/stations/route", stationHandler.GetRoute)
		}

		// Review route'larını ekle
		api.GET("/stations/:id/reviews", reviewHandler.GetStationReviews)
		api.POST("/stations/:id/reviews", reviewHandler.CreateReview)
	}

	log.Printf("Server starting on 0.0.0.0:3001")
	log.Fatal(router.Run("0.0.0.0:3001"))
}

//gurkan@Irem-MacBook-Air cmd % GOOS=linux GOARCH=amd64 go build -o charging-backend main.go
//
