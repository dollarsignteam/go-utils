package utils

import "time"

const (
	MySQLDateTimeLayout = "2006-01-02 15:04:05"
	MySQLDateLayout     = "2006-01-02"
	MySQLTimeLayout     = "15:04:05"
)

const (
	AsiaBangkokLocation  = "Asia/Bangkok"
	AsiaHongKongLocation = "Asia/Hong_Kong"
)

// Time utility instance
var Time TimeUtil

// time.Location instances
var (
	BangkokTimeLocation, _  = time.LoadLocation(AsiaBangkokLocation)
	HongKongTimeLocation, _ = time.LoadLocation(AsiaHongKongLocation)
)

// TimeUtil provides utility functions for working with time values
type TimeUtil struct{}

// ParseInBangkokLocation parses a string value
// in the Bangkok time zone with the specified layout
func (TimeUtil) ParseInBangkokLocation(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, BangkokTimeLocation)
}

// ParseInHongKongLocation parses a string value
// in the Hong Kong time zone with the specified layout
func (TimeUtil) ParseInHongKongLocation(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, HongKongTimeLocation)
}

// InBangkokTime returns a time value in the Bangkok time zone
func (TimeUtil) InBangkokTime(value time.Time) time.Time {
	return value.In(BangkokTimeLocation)
}

// InHongKongTime returns a time value in the Hong Kong time zone
func (TimeUtil) InHongKongTime(value time.Time) time.Time {
	return value.In(HongKongTimeLocation)
}

// ToMySQLDateTime returns a string value formatted
// as MySQL datetime with the specified time value
func (TimeUtil) ToMySQLDateTime(value time.Time) string {
	return value.Format(MySQLDateTimeLayout)
}

// ToMySQLDate returns a string value formatted
// as MySQL date with the specified time value
func (TimeUtil) ToMySQLDate(value time.Time) string {
	return value.Format(MySQLDateLayout)
}

// ToMySQLTime returns a string value formatted
// as MySQL time with the specified time value
func (TimeUtil) ToMySQLTime(value time.Time) string {
	return value.Format(MySQLTimeLayout)
}

// Yesterday returns a time value of yesterday
func (TimeUtil) Yesterday(value time.Time) time.Time {
	return value.AddDate(0, 0, -1)
}

// Tomorrow returns a time value of tomorrow
func (TimeUtil) Tomorrow(value time.Time) time.Time {
	return value.AddDate(0, 0, 1)
}

// IsYesterday returns true if the specified time value is yesterday
func (TimeUtil) IsYesterday(value time.Time) bool {
	return value.Format(MySQLDateLayout) == Time.Yesterday(time.Now()).Format(MySQLDateLayout)
}

// IsTomorrow returns true if the specified time value is tomorrow
func (TimeUtil) IsTomorrow(value time.Time) bool {
	return value.Format(MySQLDateLayout) == Time.Tomorrow(time.Now()).Format(MySQLDateLayout)
}

// BeginningOfDay returns a time value of the beginning of the day
func (TimeUtil) BeginningOfDay(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
}

// EndOfDay returns a time value of the end of the day
func (TimeUtil) EndOfDay(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 23, 59, 59, 999999999, value.Location())
}

// BeginningOfWeek returns a time value of the beginning of the week
func (TimeUtil) BeginningOfWeek(value time.Time) time.Time {
	daysUntilMonday := -1 * int((value.Weekday()+6)%7)
	startOfWeek := value.AddDate(0, 0, daysUntilMonday)
	return time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, value.Location())
}

// EndOfWeek returns a time value of the end of the week
func (TimeUtil) EndOfWeek(value time.Time) time.Time {
	daysUntilSunday := 7 - int(value.Weekday())
	endOfWeek := value.AddDate(0, 0, daysUntilSunday)
	return time.Date(endOfWeek.Year(), endOfWeek.Month(), endOfWeek.Day(), 23, 59, 59, 999999999, value.Location())
}

// BeginningOfMonth returns a time value of the beginning of the month
func (TimeUtil) BeginningOfMonth(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), 1, 0, 0, 0, 0, value.Location())
}

// EndOfMonth returns a time value of the end of the month
func (TimeUtil) EndOfMonth(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month()+1, 0, 23, 59, 59, 999999999, value.Location())
}

// BeginningOfYear returns a time value of the beginning of the year
func (TimeUtil) BeginningOfYear(value time.Time) time.Time {
	return time.Date(value.Year(), 1, 1, 0, 0, 0, 0, value.Location())
}

// EndOfYear returns a time value of the end of the year
func (TimeUtil) EndOfYear(value time.Time) time.Time {
	return time.Date(value.Year(), 12, 31, 23, 59, 59, 999999999, value.Location())
}
