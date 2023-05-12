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
var Time timeUtil

// time.Location instances
var (
	BangkokTimeLocation, _  = time.LoadLocation(AsiaBangkokLocation)
	HongKongTimeLocation, _ = time.LoadLocation(AsiaHongKongLocation)
)

// timeUtil provides utility functions for working with time values
type timeUtil struct{}

// ParseInBangkokLocation parses a string value
// in the Bangkok time zone with the specified layout
func (timeUtil) ParseInBangkokLocation(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, BangkokTimeLocation)
}

// ParseInHongKongLocation parses a string value
// in the Hong Kong time zone with the specified layout
func (timeUtil) ParseInHongKongLocation(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, HongKongTimeLocation)
}

// InBangkokTime returns a time value in the Bangkok time zone
func (timeUtil) InBangkokTime(value time.Time) time.Time {
	return value.In(BangkokTimeLocation)
}

// InHongKongTime returns a time value in the Hong Kong time zone
func (timeUtil) InHongKongTime(value time.Time) time.Time {
	return value.In(HongKongTimeLocation)
}

// ToMySQLDateTime returns a string value formatted
// as MySQL datetime with the specified time value
func (timeUtil) ToMySQLDateTime(value time.Time) string {
	return value.Format(MySQLDateTimeLayout)
}

// ToMySQLDate returns a string value formatted
// as MySQL date with the specified time value
func (timeUtil) ToMySQLDate(value time.Time) string {
	return value.Format(MySQLDateLayout)
}

// ToMySQLTime returns a string value formatted
// as MySQL time with the specified time value
func (timeUtil) ToMySQLTime(value time.Time) string {
	return value.Format(MySQLTimeLayout)
}
