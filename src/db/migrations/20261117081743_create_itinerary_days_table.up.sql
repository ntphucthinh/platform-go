CREATE TABLE IF NOT EXISTS itinerary_days (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    itinerary_id BIGINT NOT NULL, -- ID của itinerary cha
    day_number INT NOT NULL, -- Số thứ tự ngày
    title VARCHAR(255) NOT NULL, -- Tiêu đề chính ngày X
    subtitle VARCHAR(255) NULL DEFAULT NULL, -- Phụ đề ngắn của ngày
    summary VARCHAR(500) NULL DEFAULT NULL, -- Mô tả / Tóm tắt nội dung trải nghiệm trong ngày
    image_path VARCHAR(255) NULL DEFAULT NULL, -- Ảnh đại diện ngày
    day_tips JSON NULL DEFAULT NULL, -- Mảng chuỗi "Lưu ý trong ngày" hiển thị ở cột phải tab Ngày
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT uq_itinerary_day_number UNIQUE (itinerary_id, day_number),
    CONSTRAINT chk_itinerary_days_day_number CHECK (day_number > 0),
    CONSTRAINT fk_itinerary_days_itinerary FOREIGN KEY (itinerary_id) REFERENCES itineraries(id) ON DELETE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;
