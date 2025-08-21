package database

import (
    "context"
    "time"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "fmt"
    "os"
)

func ConnectMongoDB(uri string) (*mongo.Database, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        return nil, fmt.Errorf("MongoDB bağlantı hatası: %v", err)
    }

    // Bağlantıyı test et
    if err := client.Ping(ctx, nil); err != nil {
        return nil, fmt.Errorf("MongoDB ping hatası: %v", err)
    }

    dbName := os.Getenv("MONGODB_DB")
    if dbName == "" {
        dbName = "charging_stations"
    }

    return client.Database(dbName), nil
} 