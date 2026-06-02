CREATE TABLE IF NOT EXISTS history(
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    source VARCHAR(255),
    destination VARCHAR(255),
    original VARCHAR(255),
    translation VARCHAR(255)
);