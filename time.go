package utils

import "time"

var Time TimeUtil
var BangkokTimeLocation, _ = time.LoadLocation("Asia/Bangkok")

type TimeUtil struct{}

func (TimeUtil) ParseInBangkokLocation(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, BangkokTimeLocation)
}

func (TimeUtil) InBangkokTime(value time.Time) time.Time {
	return value.In(BangkokTimeLocation)
}

func (TimeUtil) ToMySQLDateTime(value time.Time) string {
	return value.Format("2006-01-02 15:04:05")
}

func (TimeUtil) ToMySQLDate(value time.Time) string {
	return value.Format("2006-01-02")
}

func (TimeUtil) ToMySQLTime(value time.Time) string {
	return value.Format("15:04:05")
}
