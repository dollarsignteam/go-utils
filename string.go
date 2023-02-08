package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"

	"github.com/google/uuid"
)

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

var String StringUtil

type StringUtil struct{}

func (StringUtil) RemoveDuplicateSpaces(s string) string {
	return strings.Join(strings.Fields(zeroWithReplacer.Replace(s)), " ")
}

func (StringUtil) RemoveAllSpaces(s string) string {
	return strings.Join(strings.Fields(zeroWithReplacer.Replace(s)), "")
}

func (StringUtil) UUID() string {
	return uuid.NewString()
}

func (StringUtil) MD5(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}
