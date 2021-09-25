package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFiles(t *testing.T) {
	files := main.GetFiles("./test")
	assert.Contains(t, files, "deployment.yml")
	assert.Contains(t, files, "service.yaml")
	assert.NotContains(t, files, "someotherfile.txt")
}
