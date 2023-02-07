package utils_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dollarsignteam/go-utils"
)

const (
	TestTimeLayout = "02/01/2006 15:04:05"
	TestTimeString = "07/02/2023 20:28:39"
)

var TestUTCTime, _ = time.Parse(TestTimeLayout, TestTimeString)

func TestParseInBangkokLocation(t *testing.T) {
	bangkokTime, err := utils.Time.ParseInBangkokLocation(TestTimeLayout, TestTimeString)
	expected := "2023-02-07T20:28:39+07:00"
	assert.Nil(t, err)
	assert.Equal(t, expected, bangkokTime.Format(time.RFC3339))
}

func TestInBangkokTime(t *testing.T) {
	result := utils.Time.InBangkokTime(TestUTCTime)
	expected := "2023-02-08T03:28:39+07:00"
	assert.Equal(t, expected, result.Format(time.RFC3339))
}

func TestToMySQLDateTime(t *testing.T) {
	result := utils.Time.ToMySQLDateTime(TestUTCTime)
	expected := "2023-02-07 20:28:39"
	assert.Equal(t, expected, result)
}

func TestToMySQLDate(t *testing.T) {
	result := utils.Time.ToMySQLDate(TestUTCTime)
	expected := "2023-02-07"
	assert.Equal(t, expected, result)
}

func TestToMySQLTime(t *testing.T) {
	result := utils.Time.ToMySQLTime(TestUTCTime)
	expected := "20:28:39"
	assert.Equal(t, expected, result)
}
