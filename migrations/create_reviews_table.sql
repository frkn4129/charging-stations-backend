CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    station_id VARCHAR(255) NOT NULL,
    rating DECIMAL(3,1) NOT NULL,
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_reviews_station_id ON reviews(station_id);
CREATE INDEX idx_reviews_created_at ON reviews(created_at DESC); 