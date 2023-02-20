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

var Time TimeUtil
var (
	BangkokTimeLocation, _  = time.LoadLocation(AsiaBangkokLocation)
	HongKongTimeLocation, _ = time.LoadLocation(AsiaHongKongLocation)
)

type TimeUtil struct{}

func (TimeUtil) ParseInBangkokLocation(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, BangkokTimeLocation)
}

func (TimeUtil) ParseInHongKongLocation(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, HongKongTimeLocation)
}

func (TimeUtil) InBangkokTime(value time.Time) time.Time {
	return value.In(BangkokTimeLocation)
}

func (TimeUtil) InHongKongTime(value time.Time) time.Time {
	return value.In(HongKongTimeLocation)
}

func (TimeUtil) ToMySQLDateTime(value time.Time) string {
	return value.Format(MySQLDateTimeLayout)
}

func (TimeUtil) ToMySQLDate(value time.Time) string {
	return value.Format(MySQLDateLayout)
}

func (TimeUtil) ToMySQLTime(value time.Time) string {
	return value.Format(MySQLTimeLayout)
}
