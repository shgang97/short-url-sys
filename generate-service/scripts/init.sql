CREATE DATABASE IF NOT EXISTS short_url CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE short_url;

-- 短链映射表
CREATE TABLE IF NOT EXISTS links (
    id BIGINT PRIMARY KEY,
    short_code VARCHAR(10) NOT NULL UNIQUE,
    long_url TEXT NOT NULL,
    expires_at TIMESTAMP NULL,
    click_count BIGINT UNSIGNED DEFAULT 0,
    status ENUM('active', 'disabled', 'expired') DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by VARCHAR(100),
    description VARCHAR(100),
    delete_flag varchar(1) DEFAULT 'N',
    version INT UNSIGNED DEFAULT 0,
    INDEX idx_short_code (short_code)
) COMMENT '短链映射表';

-- 缓存预热记录表
CREATE TABLE cache_warmup_logs (
    id BIGINT PRIMARY KEY,
    short_code VARCHAR(20) NOT NULL,
    warmup_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status ENUM('pending', 'success', 'failed') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by VARCHAR(100),
    description VARCHAR(100),
    delete_flag varchar(1) DEFAULT 'N',
    version INT UNSIGNED DEFAULT 0,
    INDEX idx_short_code (short_code)
) COMMENT '缓存预热记录';