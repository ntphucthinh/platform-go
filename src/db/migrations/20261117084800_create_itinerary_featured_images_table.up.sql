CREATE TABLE IF NOT EXISTS itinerary_featured_images (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    itinerary_id BIGINT NOT NULL,
    image_path VARCHAR(500) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    KEY idx_itinerary_featured_images_itinerary_id (itinerary_id),

    CONSTRAINT fk_itinerary_featured_images_itinerary
        FOREIGN KEY (itinerary_id)
        REFERENCES itineraries(id)
        ON DELETE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;
