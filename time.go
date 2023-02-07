package utils

import "time"

var Time TimeStruct
var BangkokTimeLocation, _ = time.LoadLocation("Asia/Bangkok")

type TimeStruct struct{}

func (TimeStruct) ParseInBangkokLocation(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, BangkokTimeLocation)
}

func (TimeStruct) InBangkokTime(value time.Time) time.Time {
	return value.In(BangkokTimeLocation)
}

func (TimeStruct) ToMySQLDateTime(value time.Time) string {
	return value.Format("2006-01-02 15:04:05")
}

func (TimeStruct) ToMySQLDate(value time.Time) string {
	return value.Format("2006-01-02")
}

func (TimeStruct) ToMySQLTime(value time.Time) string {
	return value.Format("15:04:05")
}
