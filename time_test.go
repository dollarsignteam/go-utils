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

func TestParseInHongKongLocation(t *testing.T) {
	hongKongTime, err := utils.Time.ParseInHongKongLocation(TestTimeLayout, TestTimeString)
	expected := "2023-02-07T20:28:39+08:00"
	assert.Nil(t, err)
	assert.Equal(t, expected, hongKongTime.Format(time.RFC3339))
}

func TestInBangkokTime(t *testing.T) {
	result := utils.Time.InBangkokTime(TestUTCTime)
	expected := "2023-02-08T03:28:39+07:00"
	assert.Equal(t, expected, result.Format(time.RFC3339))
}

func TestInHongKongTime(t *testing.T) {
	result := utils.Time.InHongKongTime(TestUTCTime)
	expected := "2023-02-08T04:28:39+08:00"
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

func TestTimeUtil(t *testing.T) {
	currentDate := time.Date(2024, 01, 31, 10, 30, 45, 123456789, time.UTC)
	tests := []struct {
		name     string
		function func(time.Time) time.Time
		input    time.Time
		expected time.Time
	}{
		{
			name:     "Yesterday",
			function: utils.Time.Yesterday,
			input:    currentDate,
			expected: time.Date(2024, 01, 30, 10, 30, 45, 123456789, time.UTC),
		},
		{
			name:     "Tomorrow",
			function: utils.Time.Tomorrow,
			input:    currentDate,
			expected: time.Date(2024, 02, 1, 10, 30, 45, 123456789, time.UTC),
		},
		{
			name:     "BeginningOfDay",
			function: utils.Time.BeginningOfDay,
			input:    currentDate,
			expected: time.Date(2024, 01, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "EndOfDay",
			function: utils.Time.EndOfDay,
			input:    currentDate,
			expected: time.Date(2024, 01, 31, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name:     "BeginningOfWeek",
			function: utils.Time.BeginningOfWeek,
			input:    currentDate,
			expected: time.Date(2024, 01, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "EndOfWeek",
			function: utils.Time.EndOfWeek,
			input:    currentDate,
			expected: time.Date(2024, 02, 4, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name:     "BeginningOfMonth",
			function: utils.Time.BeginningOfMonth,
			input:    currentDate,
			expected: time.Date(2024, 01, 01, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "EndOfMonth",
			function: utils.Time.EndOfMonth,
			input:    currentDate,
			expected: time.Date(2024, 01, 31, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name:     "BeginningOfYear",
			function: utils.Time.BeginningOfYear,
			input:    currentDate,
			expected: time.Date(2024, 01, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "EndOfYear",
			function: utils.Time.EndOfYear,
			input:    currentDate,
			expected: time.Date(2024, 12, 31, 23, 59, 59, 999999999, time.UTC),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.function(test.input)
			assert.Equal(t, test.expected, result, "they should be equal")
		})
	}
}
