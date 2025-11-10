package handler

import "time"

type TimeUnit string

const (
	Hourly  TimeUnit = "hourly"
	Daily   TimeUnit = "daily"
	Weekly  TimeUnit = "weekly"
	Monthly TimeUnit = "monthly"
)

// 获取数据库分组表达式
func (t TimeUnit) getGroupExpr() string {
	switch t {
	case Hourly:
		return "DATE_FORMAT(click_time, '%Y-%m-%d %H:00')"
	case Daily:
		return "DATE(click_time)"
	case Weekly:
		return "YEARWEEK(click_time, 1)"
	case Monthly:
		return "DATE_FORMAT(click_time, '%Y-%m')"
	default:
		return Daily.getGroupExpr()
	}
}

func (t TimeUnit) getPeriodExpr() string {
	switch t {
	case Hourly:
		return "DATE_FORMAT(click_time, '%Y-%m-%d %H:00')"
	case Daily:
		return "DATE(click_time)"
	case Weekly:
		return "YEARWEEK(click_time, 1)"
	case Monthly:
		return "DATE_FORMAT(click_time, '%Y-%m')"
	default:
		return Daily.getPeriodExpr()
	}
}

// GetDefaultDays 获取默认时间范围（天数）
func (t TimeUnit) GetDefaultDays() int {
	switch t {
	case Hourly:
		return 1 // 最近24小时
	case Daily:
		return 30 // 最近30天
	case Weekly:
		return 84 // 最近12周
	case Monthly:
		return 365 // 最近12个月
	default:
		return 30
	}
}

func (t TimeUnit) IsValid() bool {
	switch t {
	case Hourly, Daily, Weekly, Monthly:
		return true
	default:
		return false
	}
}

func (t TimeUnit) getDefaultDateRange() (*time.Time, *time.Time) {
	endDate := time.Now()
	days := t.GetDefaultDays()
	startDate := endDate.AddDate(0, 0, -days)
	return &startDate, &endDate
}
