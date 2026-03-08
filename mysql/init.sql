CREATE TABLE IF NOT EXISTS orders (
    id                 BIGINT AUTO_INCREMENT PRIMARY KEY,
    product_name       VARCHAR(255) NOT NULL,
    quantity           INT NOT NULL DEFAULT 1,
    note               TEXT,
    status             ENUM('pending', 'confirmed') NOT NULL DEFAULT 'pending',
    confirmation_token VARCHAR(64) NOT NULL,
    created_at         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
