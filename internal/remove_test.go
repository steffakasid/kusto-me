package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveFromArr(t *testing.T) {
	arr := []string{"a", "b", "c"}
	result := removeFromArray(arr, "b")
	assert.ElementsMatch(t, result, []string{"a", "c"})
}

func TestRemoveWithoutRemove(t *testing.T) {
	arr := []string{"a", "b", "c"}
	result := removeFromArray(arr, "d")
	assert.ElementsMatch(t, result, []string{"a", "b", "c"})
}
