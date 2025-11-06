package model

import "time"

// ClickEvent 点击事件表模型
type ClickEvent struct {
	ID          uint64    `gorm:"column:id;primaryKey;type:bigint unsigned" json:"id,string"`
	ShortCode   string    `gorm:"column:short_code;type:varchar(20);not null;comment:短链码" json:"short_code"`
	OriginalURL string    `gorm:"column:original_url;type:varchar(2048);not null;comment:原始URL" json:"original_url"`
	IP          string    `gorm:"column:ip;type:varchar(45);not null;comment:客户端IP" json:"ip"`
	UserAgent   string    `gorm:"column:user_agent;type:text;comment:用户代理" json:"user_agent"`
	Referer     string    `gorm:"column:referer;type:varchar(512);comment:来源页面" json:"referer"`
	Country     string    `gorm:"column:country;type:varchar(2);comment:国家代码(ISO 3166-1 alpha-2)" json:"country"`
	Region      string    `gorm:"column:region;type:varchar(100);comment:地区" json:"region"`
	City        string    `gorm:"column:city;type:varchar(100);comment:城市" json:"city"`
	DeviceType  string    `gorm:"column:device_type;type:enum('desktop','mobile','tablet','bot','other');comment:设备类型" json:"device_type"`
	Browser     string    `gorm:"column:browser;type:varchar(100);comment:浏览器" json:"browser"`
	OS          string    `gorm:"column:os;type:varchar(100);comment:操作系统" json:"os"`
	ClickTime   time.Time `gorm:"column:click_time;type:datetime(3);not null;comment:点击时间(精确到毫秒)" json:"click_time"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"created_at"`
	CreatedBy   string    `gorm:"column:created_by;type:varchar(100)" json:"created_by"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime(3)" json:"updated_at"`
	UpdatedBy   string    `gorm:"column:updated_by;type:varchar(100)" json:"updated_by"`
	Description string    `gorm:"column:description;type:varchar(100)" json:"description"`
	DeleteFlag  string    `gorm:"column:delete_flag;type:varchar(1);default:N" json:"delete_flag"`
	Version     uint      `gorm:"column:version;type:int unsigned;default:0" json:"version"`
}

// TableName 指定表名
func (ClickEvent) TableName() string {
	return "click_events"
}

// 设备类型常量
const (
	DeviceDesktop = "desktop"
	DeviceMobile  = "mobile"
	DeviceTablet  = "tablet"
	DeviceBot     = "bot"
	DeviceOther   = "other"
)
