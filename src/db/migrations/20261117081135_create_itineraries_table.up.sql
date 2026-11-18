CREATE TABLE IF NOT EXISTS itineraries (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    itinerary_name VARCHAR(255) NOT NULL, -- Tên của itinerary
    location_id BIGINT NULL DEFAULT NULL, -- ID của địa điểm
    suitable_for VARCHAR(255) NULL DEFAULT NULL, -- Phù hợp cho (VD: Cặp đôi, Gia đình...)
    description VARCHAR(500) NULL DEFAULT NULL, -- Mô tả chung về itinerary
    cover_image_path VARCHAR(255) NULL DEFAULT NULL, -- Ảnh bìa của itinerary
    travel_tips JSON NULL DEFAULT NULL, -- Mảng chuỗi "Mẹo du lịch chung" hiển thị ở tab Tổng quan
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    KEY idx_itineraries_location_id (location_id, deleted_at),

    CONSTRAINT fk_itineraries_location FOREIGN KEY (location_id) REFERENCES travel_locations(id) ON DELETE SET NULL
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;
