package utils_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dollarsignteam/go-utils"
)

var TestSpacesString = "  foo \t\r\n - bar  \u200B!  "
var TestEMVCoQRString = "00020101021229370016A000000677010111011300668776123235802TH5303764540523.99630446F5"

func TestRemoveDuplicateSpaces(t *testing.T) {
	expected := "foo - bar !"
	result := utils.String.RemoveDuplicateSpaces(TestSpacesString)
	assert.Equal(t, expected, result)
}

func TestRemoveAllSpaces(t *testing.T) {
	expected := "foo-bar!"
	result := utils.String.RemoveAllSpaces(TestSpacesString)
	assert.Equal(t, expected, result)
}

func TestUUID(t *testing.T) {
	result := utils.String.UUID()
	assert.Len(t, result, 36)
}

func TestMD5(t *testing.T) {
	result := utils.String.MD5(TestSpacesString)
	expected := "34130b8b17f2e67b2da09cd24f868885"
	assert.Equal(t, expected, result)
}

func TestSHA1(t *testing.T) {
	result := utils.String.SHA1(TestSpacesString)
	expected := "5045a76071cb10b5fa05b73af1e8e174d4979955"
	assert.Equal(t, expected, result)
}

func TestSHA256(t *testing.T) {
	result := utils.String.SHA256(TestSpacesString)
	expected := "bafa334ba4639eca91f087ad98a0dcc9d1ac2f82da8beafed8fbaad717a51c6d"
	assert.Equal(t, expected, result)
}

func TestHashAndVerifyPassword(t *testing.T) {
	password := "foo"
	result, err := utils.String.HashPassword(password)
	assert.Nil(t, err)
	err = utils.String.VerifyPassword(result, password)
	assert.Nil(t, err)
}

func TestHashPassword_Error(t *testing.T) {
	password := make([]byte, 80)
	result, err := utils.String.HashPassword(string(password))
	assert.Empty(t, result)
	assert.EqualError(t, err, "bcrypt: password length exceeds 72 bytes")
}

func TestHashCrc32(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "Test Case 1", input: "Hello, World!", expected: "ec4ac3d0"},
		{name: "Test Case 2", input: "Lorem Ipsum", expected: "358ad45d"},
		{name: "Test Case 3", input: "1234567890", expected: "261daee5"},
		{name: "Test Case 4", input: "A4D7B2B7-D62D-423C-B0C2-2A871F98E427", expected: "042b1405"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.String.HashCrc32(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestParseEMVCoQRString(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		errExpected  error
		infoExpected utils.EMVCoQRInfo
	}{
		{
			name:        "Test Case 1",
			input:       "",
			errExpected: fmt.Errorf("invalid specified qr string length"),
		},
		{
			name:        "Test Case 2",
			input:       "00020101021229370016A000000677010111011300668776123235801TH5303764540523.99630446F5",
			errExpected: fmt.Errorf("qr checksum is incorrect"),
		},
		{
			name:  "Test Case 3",
			input: "00020101021229370016A000000677010111011300668776123235802TH5303764540523.99630446F5",
			infoExpected: utils.EMVCoQRInfo{
				Format:          "12",
				MerchantAccount: "0016A00000067701011101130066877612323",
				Amount:          "23.99",
				PhoneNumber:     "66877612323",
				CountryCode:     "TH",
				Crc:             "46F5",
				CurrencyISO4217: "764",
			},
		},
		{
			name:  "Test Case 4",
			input: "00020101021230730016A00000067701011201150107536000315080214KB0000018898870312KPSX8YYB3JO853037645406100.005802TH62130709X8YYB3JO863049EEA",
			infoExpected: utils.EMVCoQRInfo{
				Format:          "12",
				MerchantAccount: "0016A00000067701011201150107536000315080214KB0000018898870312KPSX8YYB3JO8",
				Amount:          "100.00",
				PhoneNumber:     "",
				CountryCode:     "TH",
				Crc:             "9EEA",
				CurrencyISO4217: "764",
				BillerID:        "010753600031508",
				Ref1:            "KB000001889887",
				Ref2:            "KPSX8YYB3JO8",
				Ref3:            "X8YYB3JO8",
			},
		},
		{
			name:  "Test Case 5",
			input: "00020101021253037645802TH29370016A000000677010111021303455660038075406200.006304ABA1",
			infoExpected: utils.EMVCoQRInfo{
				Format:          "12",
				MerchantAccount: "0016A00000067701011102130345566003807",
				Amount:          "200.00",
				PhoneNumber:     "0345566003807",
				CountryCode:     "TH",
				Crc:             "ABA1",
				CurrencyISO4217: "764",
				BillerID:        "",
				Ref1:            "",
				Ref2:            "",
				Ref3:            "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			info, err := utils.String.ParseEMVCoQRString(test.input)
			assert.Equal(t, test.errExpected, err)
			assert.Equal(t, test.infoExpected, info)
		})
	}
}

func TestAESEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		plaintext string
		expected  string
	}{
		{
			name:      "empty string",
			key:       "ExampleKey123456",
			plaintext: "",
			expected:  "",
		},
		{
			name:      "normal string",
			key:       "ExampleKey123456",
			plaintext: "Hello, World!",
			expected:  "Hello, World!",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cipherText, err := utils.String.AESEncrypt(test.key, test.plaintext)
			assert.NoError(t, err)
			decryptedText, err := utils.String.AESDecrypt(test.key, cipherText)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, decryptedText)
		})
	}
}

func BenchmarkRemoveDuplicateSpaces(b *testing.B) {
	for n := 0; n < b.N; n++ {
		utils.String.RemoveDuplicateSpaces(TestSpacesString)
	}
}

func BenchmarkRemoveAllSpaces(b *testing.B) {
	for n := 0; n < b.N; n++ {
		utils.String.RemoveAllSpaces(TestSpacesString)
	}
}

func BenchmarkMD5(b *testing.B) {
	for n := 0; n < b.N; n++ {
		utils.String.MD5(TestSpacesString)
	}
}

func BenchmarkSHA1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		utils.String.SHA1(TestSpacesString)
	}
}

func BenchmarkSHA256(b *testing.B) {
	for n := 0; n < b.N; n++ {
		utils.String.SHA256(TestSpacesString)
	}
}

func BenchmarkHashPassword(b *testing.B) {
	for n := 0; n < b.N; n++ {
		utils.String.HashPassword(TestSpacesString)
	}
}

func BenchmarkHashCrc32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utils.String.HashCrc32(TestSpacesString)
	}
}

func BenchmarkParseEMVCoQRString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utils.String.ParseEMVCoQRString(TestEMVCoQRString)
	}
}
