CREATE DATABASE IF NOT EXISTS short_url CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE short_url;

-- 点击事件明细表
CREATE TABLE IF NOT EXISTS click_events (
    id BIGINT PRIMARY KEY,
    short_code VARCHAR(20) NOT NULL COMMENT '短链码',
    original_url VARCHAR(2048) NOT NULL COMMENT '原始URL',
    ip VARCHAR(45) NOT NULL COMMENT '客户端IP',
    user_agent TEXT COMMENT '用户代理',
    referer VARCHAR(512) COMMENT '来源页面',
    country VARCHAR(2) COMMENT '国家代码(ISO 3166-1 alpha-2)';
    region VARCHAR(100) COMMENT '地区',
    city VARCHAR(100) COMMENT '城市',
    device_type ENUM('desktop', 'mobile', 'tablet', 'bot', 'other') COMMENT '设备类型',
    browser VARCHAR(100) COMMENT '浏览器',
    os VARCHAR(100) COMMENT '操作系统',
    click_time DATETIME(3) NOT NULL COMMENT '点击时间(精确到毫秒)',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    created_by VARCHAR(100),
    updated_at DATETIME(3),
    updated_by VARCHAR(100),
    description VARCHAR(100),
    delete_flag varchar(1) DEFAULT 'N',
    version INT UNSIGNED DEFAULT 0,
    INDEX idx_short_code_time (short_code, click_time),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='点击事件明细表';