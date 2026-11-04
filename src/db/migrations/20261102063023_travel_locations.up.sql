CREATE TABLE IF NOT EXISTS travel_locations (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    country_code CHAR(2) NOT NULL,
    country_name VARCHAR(255) NOT NULL,
    state_code VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    latitude DECIMAL(10, 7) NULL,
    longitude DECIMAL(10, 7) NULL,
    image_path VARCHAR(255) NULL,
    status ENUM ('visited', 'wishlist') NOT NULL DEFAULT 'visited',
    description TEXT NULL,
    notes TEXT NULL,
    google_maps_url VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    -- Đảm bảo mỗi tỉnh/bang của một quốc gia chỉ xuất hiện duy nhất 1 lần trên bản đồ khi chưa xóa
    UNIQUE KEY uq_travel_locations_country_state (country_code, state_code, (IF(deleted_at IS NULL, 1, NULL))),

    KEY idx_travel_locations_country (country_code),
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;
