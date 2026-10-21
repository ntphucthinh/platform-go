CREATE TABLE IF NOT EXISTS homestay_images (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    homestay_id BIGINT NOT NULL,
    image_path TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    room_id BIGINT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    is_cover BOOLEAN NOT NULL DEFAULT FALSE,

    KEY idx_homestay_images_homestay_id (homestay_id),
    KEY idx_homestay_images_room_id (room_id),

    CONSTRAINT fk_homestay_images_homestay
        FOREIGN KEY (homestay_id)
        REFERENCES homestays(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_homestay_images_room
        FOREIGN KEY (room_id)
        REFERENCES homestay_rooms(id)
        ON DELETE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;