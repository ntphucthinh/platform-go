CREATE TABLE IF NOT EXISTS itinerary_themes (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    itinerary_id BIGINT NOT NULL, -- ID của itinerary
    theme_id BIGINT NOT NULL, -- ID của travel_themes
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE KEY uq_itinerary_themes (itinerary_id, theme_id),

    KEY idx_itinerary_themes_theme_id (theme_id),

    CONSTRAINT fk_itinerary_themes_itinerary FOREIGN KEY (itinerary_id) REFERENCES itineraries(id) ON DELETE CASCADE,
    CONSTRAINT fk_itinerary_themes_theme FOREIGN KEY (theme_id) REFERENCES travel_themes(id) ON DELETE RESTRICT
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;
