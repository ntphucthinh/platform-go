CREATE TABLE IF NOT EXISTS homestay_rooms (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    homestay_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NULL,
    price_amount DECIMAL(15, 2) NULL,
    price_unit VARCHAR(50) NULL,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    KEY idx_homestay_rooms_homestay_id (homestay_id),

    CONSTRAINT fk_homestay_rooms_homestay
        FOREIGN KEY (homestay_id)
        REFERENCES homestays(id)
        ON DELETE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;