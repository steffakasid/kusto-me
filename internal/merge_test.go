package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeMaps(t *testing.T) {
	map1 := map[string]string{"key1": "value1"}
	map2 := map[string]string{"key2": "value2"}
	map3 := mergeMaps(map1, map2)
	assert.Contains(t, map3, "key1")
	assert.Contains(t, map3, "key2")
}

func TestMergeArrays(t *testing.T) {
	arr1 := []string{"entry1"}
	arr2 := []string{"entry2"}
	arr3 := mergeArrays(arr1, arr2)
	assert.ElementsMatch(t, []string{"entry1", "entry2"}, arr3)
}
