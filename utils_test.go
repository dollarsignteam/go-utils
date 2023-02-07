package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dollarsignteam/go-utils"
)

func TestPointerOf(t *testing.T) {
	value := 0
	result := utils.PointerOf(value)
	assert.Equal(t, &value, result)
}

func TestPackageName(t *testing.T) {
	name := utils.PackageName()
	assert.Equal(t, "go-utils_test", name)
}
