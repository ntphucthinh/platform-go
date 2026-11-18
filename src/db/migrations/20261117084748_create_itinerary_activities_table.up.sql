CREATE TABLE IF NOT EXISTS itinerary_activities (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    itinerary_day_id BIGINT NOT NULL,
    start_time TIME NOT NULL, -- Giờ bắt đầu định dạng HH:mm (VD: "07:30", "14:00")
    title VARCHAR(255) NOT NULL, -- Tên hoạt động
    activity_category ENUM('dining', 'sightseeing', 'nature', 'cafe', 'shopping', 'rest', 'other') NULL DEFAULT NULL, -- 'Ăn uống' | 'Tham quan' | 'Thiên nhiên' | 'Cà phê' | 'Mua sắm' | 'Nghỉ ngơi'
    duration_minutes SMALLINT UNSIGNED NULL,  -- Thời lượng tính bằng phút
    location VARCHAR(500) NULL, -- Địa điểm cụ thể 
    google_maps_url VARCHAR(500) NULL, -- Link Google Maps
    description VARCHAR(500) NULL, -- Mô tả hoạt động
    notes VARCHAR(255) NULL, -- Ghi chú thêm
    image_path VARCHAR(255) NULL, -- Ảnh đại diện hoạt động 
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    KEY idx_itinerary_activities_day_time (itinerary_day_id, start_time, id),

    CONSTRAINT fk_itinerary_activities_day FOREIGN KEY (itinerary_day_id) REFERENCES itinerary_days(id) ON DELETE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;
