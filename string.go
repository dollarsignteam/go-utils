package utils

import (
	"crypto/md5"  // #nosec
	"crypto/sha1" // #nosec
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/google/uuid"
)

// Constants representing various Unicode characters
const (
	zeroWidthSpace        = '\u200B'
	zeroWidthNoBreakSpace = '\uFEFF'
	wordJoiner            = '\u2060'
	zeroWidthJoiner       = '\u200D'
	leftToRightMark       = '\u200E'
	rightToLeftMark       = '\u200F'
	noBreakingSpace       = '\u00A0'
	empty                 = ""
)

var zeroWithReplacer = strings.NewReplacer(
	string(zeroWidthSpace), empty,
	string(zeroWidthNoBreakSpace), empty,
	string(wordJoiner), empty,
	string(zeroWidthJoiner), empty,
	string(leftToRightMark), empty,
	string(rightToLeftMark), empty,
	string(noBreakingSpace), empty,
)

// String utility instance
var String StringUtil

// StringUtil provides utility functions for manipulating strings
type StringUtil struct{}

// RemoveDuplicateSpaces removes duplicate spaces from the input string
func (StringUtil) RemoveDuplicateSpaces(s string) string {
	return strings.Join(strings.Fields(zeroWithReplacer.Replace(s)), " ")
}

// RemoveAllSpaces removes all spaces from the input string
func (StringUtil) RemoveAllSpaces(s string) string {
	return strings.Join(strings.Fields(zeroWithReplacer.Replace(s)), "")
}

// UUID generates a new UUID string
func (StringUtil) UUID() string {
	return uuid.NewString()
}

// MD5 generates an MD5 hash for the input string
func (StringUtil) MD5(s string) string {
	hash := md5.Sum([]byte(s)) // #nosec
	return hex.EncodeToString(hash[:])
}

// SHA1 generates a SHA1 hash for the input string
func (StringUtil) SHA1(s string) string {
	hash := sha1.Sum([]byte(s)) // #nosec
	return hex.EncodeToString(hash[:])
}

// SHA256 generates a SHA256 hash for the input string
func (StringUtil) SHA256(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}
