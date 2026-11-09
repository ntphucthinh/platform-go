CREATE TABLE IF NOT EXISTS travel_memory_photos (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    memory_day_id BIGINT NOT NULL,
    image_path VARCHAR(255) NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT chk_travel_memory_photos_sort_order CHECK (sort_order >= 0),

    KEY idx_travel_memory_photos_day_id (memory_day_id)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;
